#include "DatxRecharge.hpp"
#include <DatxosLib/multi_index.hpp>
#include <string>
#include <DatxosLib/print.hpp>
#include <DatxosLib/crypto.h>
#include <DatxosLib/chain.h>

namespace datxos
{
/// @abi action
void recharge::recorduser(account_name datxaddress,string address,account_name bpname)
{
    require_auth(bpname);
    account_name producers[21]; 
    uint32_t bytes_populated = get_active_producers(producers, sizeof(account_name)*21); 
        bool Isproducer = false; 
        for (int i = 0; i < sizeof(producers)/sizeof(account_name) ;i++){ 
             if(producers[i] == bpname) 
             Isproducer = true; 
        } 
        datxos_assert(Isproducer, "this func can only be called by producers");

        char *data=(char*)address.c_str();
        checksum256 adrhash;
        sha256( data, sizeof(data), &adrhash );

        print("recorduser:  ",N(adrhash), "\n");

    // define the table 
        users users_table(_self,_self);
        auto idx = users_table.template get_index<N(hash)>();
        auto idu = idx.find(get_hash(adrhash));
        datxos_assert(idu == idx.end(), "This address is already exists");

        users_table.emplace(_self, [&](auto &u) {
            u.id= users_table.available_primary_key();
            u.hash = adrhash;
            u.datxaddress=datxaddress;
            u.address=address;
        });    
}
/// @abi action
void recharge::charge(account_name bpname,string hash,string from,string to,int64_t blocknum,string quantity,string category,string memo)
  {
        require_auth(bpname);
    //check the caller is producer
        account_name producers[21]; 
        uint32_t bytes_populated = get_active_producers(producers, sizeof(account_name)*21); 
        bool Isproducer = false; 
        for (int i = 0; i < sizeof(producers)/sizeof(account_name) ;i++){ 
             if(producers[i] == bpname) 
             Isproducer = true; 
        } 
        datxos_assert(Isproducer, "this func can only be called by producers");
        
    // get all expired transactions and remove to expirdrecord
       this->expired_trx();
   // check the countrecord whether exist
        countrecords count_table(_self,_self);
        perfections success_table(_self,_self);
        expirations expired_table(_self,_self);
        records records_table(_self,_self);

    // caculate hash for trx data

        string str=hash+from+to+std::to_string(blocknum)+quantity+category+memo;
        char *data=(char*)str.c_str();
        checksum256 calc_hash;
        sha256( data, sizeof(data), &calc_hash );

    //check if the hash alerday success
        auto idx2 = success_table.template get_index<N(data)>();
        auto itr2 = idx2.find(get_hash(calc_hash) );
        datxos_assert(itr2==idx2.end(), "This transaction is already success");

    //check if the hash already push more than three times
        auto idx3 = count_table.template get_index<N(data)>();
        auto itr3 = idx3.find(get_hash(calc_hash));
        if (itr3 != idx3.end()){
            datxos_assert(itr3->count < 4, "This transaction is already expire more than three times");
        }       

    //check the data hash whether equal
        auto idx = records_table.template get_index<N(data)>();
        auto itr = idx.find(get_hash(calc_hash));
        if (itr == idx.end()){
            // if data hash isn't exist then record it in record table 
            vector<account_name> producers_name;
            producers_name.push_back(bpname);
            records_table.emplace(_self, [&](auto &s) {
                s.id= records_table.available_primary_key();
                s.trxid= hash;
                s.from = from;
                s.to = to;
                s.blocknum= blocknum;
                s.quantity=quantity;
                s.category=category;
                s.memo = memo;
                s.data=calc_hash; 
                s.start_time=now(); //seconds
                s.producers=producers_name;
            });  
            return;
        }

        //check producer whether more than 15
        int size =(int)itr->producers.size();
        if(size<5)
        {
            auto itr3 = std::find( itr->producers.cbegin(), itr->producers.cend(), bpname );
            datxos_assert(itr3 == itr->producers.end(), "This producer already initiated a request for this transaction!");  
            records_table.modify(*itr, get_self(), [&](auto& p){
                    p.producers.push_back(bpname);
                });
        }else
        {
        //transanfer to datx account 
        // you should add DatxRecharge permission first, and then can send action DatxToken contract.

            int64_t q=atoll(quantity.c_str());
            if(category=="EOS"){
                //eos "to" is memo
                asset quant{q,S(4,DEOS)};
                action(
                    permission_level{N(datxos.deos), N(active) },
                    N(datxos.dtoke), N(transfer),
                    std::make_tuple(N(datxos.deos), N(memo), quant,memo)
                ).send();
            }else{
                char *data=(char*)from.c_str();
                checksum256 adrhash;
                sha256( data, sizeof(data), &adrhash );

                print("charge:  ",N(adrhash), "\n");


                users users_table(_self,_self);
                auto idu = users_table.template get_index<N(hash)>();
                auto idxu = idu.find(get_hash(adrhash));
                datxos_assert(idxu != idu.end(), "This address is not exists");
                //get the user's datxaddress by btc's or eth's address
                if(category=="BTC"){
                    asset quant{q,S(4,DBTC)};
                    action(
                        permission_level{N(datxos.dbtc), N(active) },
                        N(datxos.dtoke), N(transfer),
                        std::make_tuple(N(datxos.dbtc), N(idxu->datxaddress), quant,memo)
                    ).send();
                }else if(category=="ETH"){
                    asset quant{q,S(4,DETH)};
                    action(
                        permission_level{N(datxos.deth), N(active) },
                        N(datxos.dtoke), N(transfer),
                        std::make_tuple(N(datxos.deth), N(idxu->datxaddress), quant,memo)
                    ).send();
                }else{
                    datxos_assert(false, "This category is not support");
                }            
            }

            //after transfer to datxaddress remove info to success
            success_table.emplace(_self, [&](auto &s) 
            {
                s.id= success_table.available_primary_key();
                s.trxid= hash;
                s.from = from;
                s.to = to;
                s.blocknum= blocknum;
                s.quantity=quantity;
                s.category=category;
                s.memo = memo;
                s.data=calc_hash; 
            });   

            //remove
            records_table.erase(*itr);//remove the transaction in the record_table
            auto in = expired_table.template get_index<N(data)>();
            auto r = in.find(get_hash(calc_hash) );
            if(r!=in.end())
            {
                expired_table.erase(*r);
            }

            auto idx4 = count_table.template get_index<N(data)>();
            auto itr4 = idx4.find(get_hash(calc_hash));
            if (itr4 != idx4.end()){
                count_table.erase(*itr4);
            } 
        }
}

//get the expired transaction from record
void recharge::expired_trx()
{
    records trans_table(_self,_self);
    int count = 0;
    for ( auto it = trans_table.begin(); it != trans_table.end() && count<100;)
    { 
        countrecords count_table(_self,_self);
        uint64_t subtime = now() - it->start_time;
        if(it->producers.size() < 15&&subtime > 5*63)
        {
            expirations expire_table(_self,_self);
            auto idx = expire_table.template get_index<N(data)>();
            auto idf = idx.find(get_hash(it->data));
            if (idf == idx.end()){
                expire_table.emplace(_self, [&](auto &e) {
                    e.id= expire_table.available_primary_key();
                    e.trxid=it->trxid;
                    e.from=it->from;
                    e.blocknum=it->blocknum;
                    e.quantity=it->quantity;
                    e.category=it->category;
                    e.memo=it->memo;
                    e.data=it->data;        
                }); 
                it = trans_table.erase(it);
            }else{
                auto idc =count_table.template get_index<N(data)>();
                auto idcc = idc.find(get_hash(it->data));
                if (idcc != idc.end()){
                    count_table.modify(*idcc, get_self(), [&](auto& c)
                    {
                        c.count=idcc->count+1;
                    });
                }else{
                   count_table.emplace(_self, [&](auto &co){
                        co.id=count_table.available_primary_key();
                        co.data=it->data;
                        co.count=1;
                    });
                }
                ++it;
            }
            ++count;
        }else
        {
            ++it;
        }
    }
}

} // namespace Datxos


// modify form DATXOS_ABI
#define DATXOS_ABI_EX( TYPE, MEMBERS ) \
extern "C" { \
   void apply( uint64_t receiver, uint64_t code, uint64_t action ) { \
      auto self = receiver; \
      if( action == N(onerror)) { \
         /* onerror is only valid if it is for the "Datxos" code account and authorized by "Datxos"'s "active permission */ \
         datxos_assert(code == N(Datxos), "onerror action's are only valid from the \"Datxos\" system account"); \
      } \
      if( code == self || code == N(datxos.dtoke) || action == N(onerror) ) { \
         TYPE thiscontract( self ); \
         switch( action ) { \
            DATXOS_API( TYPE, MEMBERS ) \
         } \
         /* does not allow destructor of thiscontract to run: Datxos_exit(0); */ \
      } \
   } \
}

DATXOS_ABI_EX(datxos::recharge,(recorduser)(charge))
