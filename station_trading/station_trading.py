import ctypes
from io import StringIO
import pandas as pd
import csv
import requests
import bz2
import os.path
import pandas as pd
import time
import numpy as np

pd.set_option('display.float_format', lambda x: '%.1f' % x)

def load_static_dump(url: str, overwrite: bool = False) -> pd.DataFrame:
    """
    Downloads a file from the fuzzworks static dump in bz2 format,
    saves it as a csv and loads it into a pandas dataframe
    """
    filename = url.split('/')[-1][:-4]

    if not os.path.isfile('SDE/'+filename) or overwrite:

        response = requests.get(url)
        text = bz2.decompress(response.content)

        with open('SDE/'+filename, 'w', encoding = 'utf-8') as file:
            for line in StringIO(text.decode()):
                file.write(line)

    return pd.read_csv('SDE/'+filename)

def deserialize_marketGroupID(groupID: int) -> list:
    """
    Due to the recursive nature of market groups, not all IDs in invMarketGroups have items
    This function deserializes a marketGroupID into a list of marketGroupIDs that actually have an item
    """
    tmp = invMarketGroups.loc[invMarketGroups['parentGroupID'] == str(groupID)]
    types = list(invMarketGroups.loc[(invMarketGroups['hasTypes'] == 1)&(invMarketGroups['marketGroupID'] == groupID)]['marketGroupID'])

    hasTypes = []

    while not hasTypes or list(tmp):

        hasTypes = list(tmp.loc[tmp['hasTypes'] == 1]['marketGroupID'].astype(str))
        types.extend(hasTypes)
        tmp = invMarketGroups[invMarketGroups['parentGroupID'].isin(list(tmp['marketGroupID'].astype(str)))]
        
        if tmp.empty and not hasTypes:
            return types
        
    return types

invTypes = load_static_dump('https://www.fuzzwork.co.uk/dump/latest/invTypes.csv.bz2')
invMetaTypes = load_static_dump('https://www.fuzzwork.co.uk/dump/latest/invMetaTypes.csv.bz2')
invMarketGroups = load_static_dump('https://www.fuzzwork.co.uk/dump/latest/invMarketGroups.csv.bz2')

length = len(invTypes)

# Remove items that aren't published and can't be traded on the market

invTypes = invTypes.loc[(invTypes['published'] == 1)&(invTypes['marketGroupID'] != 'None')]

print(f"[+] Removed {length-len(invTypes)} items that aren't published or can't be traded on the market")
length = len(invTypes)

# Remove Apparel, Blueprints, Skillbooks, SKINS, Pilot's Services, Special Edition Assets, Materials, Trade Goods, Ships, Rigs

deserialized = [] # contains groupIDs that have items
groupIDs = [2, 1396, 1954, 150, 1922, 1659, 1031, 19, 4, 475, 1111]  # marketGroupIds to be excluded

for id in groupIDs:
    deserialized.extend(deserialize_marketGroupID(id))

invTypes = invTypes[~invTypes['marketGroupID'].isin(deserialized)]

print(f"[+] Removed {length-len(invTypes)} items that are in one of the following groups: Apparel, Blueprints, Skillbooks, SKINS, Pilot's Services, Special Edition Assets, Materials, Trade Goods, Ships, Rigs")
length = len(invTypes)

# Join with invMetaGroups to get the metaGroupID for each item

invTypes = pd.merge(invTypes, invMetaTypes, how="left", on= "typeID").drop("parentTypeID", axis = 1)

# Remove all T2 and officer items

invTypes = invTypes[~invTypes['metaGroupID'].isin([2, 5])]

print(f"[+] Removed {length-len(invTypes)} items that are T2 or Officer\n")

# Export the (filtered) typeIDs to a .txt file so it can be read later and remove the last line

invTypes['typeID'].to_csv('Shared/items.txt', index = False, header = False, line_terminator='\n')

# Read the exported function from the compiled .dll

library = ctypes.cdll.LoadLibrary('Libraries/esiRequests.dll')

pull_market_data = library.pullMarketData
pull_market_data.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int, ctypes.c_int]  #inputfile, outputfile, concurrency, region_id

pull_market_data("Shared/items.txt".encode("utf-8"), "Shared/prices.json".encode("utf-8"), 50, 10000002)

# Load the market data into a dataframe

prices = pd.read_json("Shared/prices.json")

# Merge with invTypes to get each item's typeName

prices = pd.merge(prices, invTypes[['typeID', 'typeName']], how = "left", left_on = "TypeID", right_on = "typeID").drop('typeID', axis = 1)
prices = prices.rename(columns = {"typeName": "TypeName"})

# Pull all market orders in the TTT

REFRESH_TOKEN = "" # your refresh token here
CLIENT_ID = "" # your client id here


pull_structure_orders = library.pullStructureOrders
pull_structure_orders.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_longlong, ctypes.c_int, ctypes.c_char_p] # refresh_token, client_id, structure_id, concurrency, outputfile

pull_structure_orders(REFRESH_TOKEN.encode('utf-8'), CLIENT_ID.encode('utf-8'), 1028858195912, 20,"Shared/tttOrders.json".encode('utf-8'))

# Update the Max Buy column with the TTT buy orders

tttOrders = pd.read_json('Shared/tttOrders.json')
tttOrders = tttOrders.loc[tttOrders['is_buy_order']==True] # remove sell orders
tttOrders['tttBought'] = tttOrders['volume_total'] - tttOrders['volume_remain']

tttOrders = tttOrders.groupby('type_id').agg({'price': 'max', 'tttBought': 'sum'}).reset_index() # Get the max buy value and the amount of items that have been bought through those orders

prices = pd.merge(prices, tttOrders, how="left", left_on="TypeID", right_on = "type_id").drop("type_id", axis = 1).rename(columns = {'price': "tttMaxBuy"}) # add them to the dataframe with all the prices
prices['MaxBuy'] = prices[['MaxBuy', 'tttMaxBuy']].max(axis = 1)
prices['Bought'] = prices['Bought'] + prices['tttBought']

# remove items with no buy orders
prices = prices.loc[prices['MaxBuy'] != 0]

# Calculate margin and isk volumes

prices['Margin'] = prices['MinSell'] / prices['MaxBuy']
prices['ISKVolume'] = (prices['MaxBuy']+prices['MinSell'])/2 * prices['AvgVolume']
prices['SoldVolume'] = prices['Sold'] * prices['MinSell']
prices['BoughtVolume'] = prices['Bought'] * prices['MaxBuy']
prices['SellToBuy'] = prices['Sold']/prices['Bought']

# Remove all items with too little margin and too little volume and which haven't been bought or sold

prices = prices.loc[(prices['Margin'] > 1.3)&(prices['Margin'] < 5)]
prices = prices.loc[prices['ISKVolume'] > 300000000]
prices = prices.loc[prices['AvgVolume'] > 1]

prices.replace([np.inf, -np.inf], np.nan, inplace=True)
prices.fillna(0, inplace = True)

prices = prices.loc[prices['SellToBuy'] > 0]

# Create a measure based on which the items will be sorted
# (ISKVolume + SoldVolume + BoughtVolume) / ((BuyQuant/BuyOrders)*(SellQuant/SellOrders))

prices['Sorting'] = (prices['ISKVolume'] + prices['SoldVolume'] + prices['BoughtVolume']) / ((prices['BuyQuant']/prices['BuyOrders'])*(prices['SellQuant']*prices['SellOrders']))

print(prices.sort_values('Sorting', ascending = False))