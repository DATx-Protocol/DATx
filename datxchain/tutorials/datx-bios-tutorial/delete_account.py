import json
with open('accounts.json','r') as f:
        a = json.load(f)
        del a['users'][10:]  # 4
        del a['producers'][5:]  # 5

with open('new_accounts.json','w') as f:
        json.dump(a,f)
