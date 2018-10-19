#include "DatxExtract.hpp"
#include <DatxosLib/multi_index.hpp>
#include <DatxosLib/chain.h>

namespace datxos
{
    /// @abi action
    void extract::recordtrx(transaction_id_type trxid, account_name producer,string category)
    {   
        //handle expire transaction and rollback transaction
        this->updateexpire();
        this->rollbacktrx();

        require_auth(producer);
        account_name producers[21]; 
        uint32_t bytes_populated = get_active_producers(producers, sizeof(account_name)*21); 
        bool Isproducer = false; 
        for (int i = 0; i < sizeof(producers)/sizeof(account_name) ;i++){ 
            if(producers[i] == producer) 
            Isproducer = true; 
            break;
        } 
        datxos_assert(Isproducer, "this func can only be called by producers");

        transrecords transrecord_table(N(Datxos.token),N(Datxos.token));
        auto idx0 = transrecord_table.template get_index<N(fixed_key)>();
        auto itr0 = idx0.find( get_fixed_key(trxid) );
        datxos_assert(itr0 != idx0.end(), "trxid not exists in token transfer");

        datxos_assert(itr0 ->quantity.symbol.name() == string_to_name(category.c_str())
        ,"category not correct");


        records trans_table(_self,_self);
        auto idx = trans_table.template get_index<N(fixed_key)>();
        auto itr = idx.find( get_fixed_key(trxid) );
        datxos_assert(itr==idx.end(), "trxid already exists");

        successtrxs success_table(_self,_self);
        auto idx2 = success_table.template get_index<N(fixed_key)>();
        auto itr2 = idx2.find(get_fixed_key(trxid) );
        datxos_assert(itr2==idx2.end(), "trxid already success");

        expirations expire_table(_self,_self);
        auto idx3 = expire_table.template get_index<N(fixed_key)>();
        auto itr3 = idx3.find(get_fixed_key(trxid) );
        if(itr3 != idx3.end()){
            expire_table.erase(*itr3);
        }

        trans_table.emplace(_self, [&](auto &s) {
            s.id= trans_table.available_primary_key();
            s.trxid= trxid;
            s.start_time = now();
            s.producer = producer;
            s.category = category;
        }); 

    }


    /// @abi action
    void extract::setverifiers(vector<account_name> accounts){
        //handle expire transaction and rollback transaction
        this->updateexpire();
        this->rollbacktrx();

        require_auth(_self);

        verifiers veri_table(_self,_self);
        for(auto itr = veri_table.begin(); itr != veri_table.end();) {
            itr = veri_table.erase(itr);
        }

        for( const auto& a : accounts ){
            veri_table.emplace(_self, [&](auto &s) {
                s.account = a;
            });
        }
    }

    /// @abi action
    void extract::setdoing(transaction_id_type trxid, account_name producer,account_name verifier){
        //handle expire transaction and rollback transaction
        this->updateexpire();
        this->rollbacktrx();

        // require_auth(verifier);
        verifiers veri_table(_self,_self);
        auto vidx = veri_table.find(verifier);
        datxos_assert(vidx != veri_table.end(),"verifier is not exists");

        records trans_table(_self,_self);
        auto idx = trans_table.template get_index<N(fixed_key)>();
        auto itr = idx.find(get_fixed_key(trxid) );
        datxos_assert(itr != idx.end(), "trxid not exists");
        datxos_assert(itr->producer == producer, "trxid not this producer");
        
        auto itr2 = std::find( itr->verifiers.begin(), itr->verifiers.end(), verifier );
        print("verifier name",N(*itr2),"\n");
        datxos_assert( itr2 == itr->verifiers.end(), "this verifier has comfirmed this trxid" );

        trans_table.modify(*itr, get_self(), [&](auto& p)
                                              {
                                                p.verifiers.push_back(verifier);
                                              });
        if(itr->verifiers.size() > 15){
            trans_table.modify(*itr, get_self(), [&](auto& p)
                                              {
                                                p.countdown_time = now();
                                              });
        }
    }
    
    /// @abi action
    void extract::setsuccess(transaction_id_type trxid ,account_name producer){
        //handle expire transaction and rollback transaction
        this->updateexpire();
        this->rollbacktrx();

        require_auth(producer);

        account_name producers[21]; 
        uint32_t bytes_populated = get_active_producers(producers, sizeof(account_name)*21); 
        bool Isproducer = false; 
        for (int i = 0; i < sizeof(producers)/sizeof(account_name) ;i++){ 
             if(producers[i] == producer) 
             Isproducer = true; 
             break;
        } 
        datxos_assert(Isproducer, "this func can only be called by producers");

        records trans_table(_self,_self);
        auto idx = trans_table.template get_index<N(fixed_key)>();
        auto itr = idx.find(get_fixed_key(trxid) );
        datxos_assert(itr != idx.end(), "trxid not in doing records");

        trans_table.modify(*itr, get_self(), [&](auto& p)
                                              {
                                                p.successconfirm.push_back(producer);
                                              });
        
        if (itr -> successconfirm.size() > 15) {
            successtrxs success_table(_self,_self);
            success_table.emplace(_self, [&](auto &s) {
                s.id= success_table.available_primary_key();
                s.trxid= trxid;
                s.producer = itr -> producer;
                s.timestamp = now();
                s.category = itr->category;
            });
            trans_table.erase(*itr);
        }
        
    }

    /// @abi action
    void extract::updateexpire(){
        records trans_table(_self,_self);
        auto idx = trans_table.template get_index<N(start_time)>();
        int count = 0;
        for ( auto it = idx.cbegin(); it != idx.cend() && count < 100;) {
            uint64_t subtime = now() - it->start_time;
            if(subtime <= 5 * 60){
                break;
            }
            if(it->verifiers.size() < 15){
                expirations expire_table(_self,_self);
                expire_table.emplace(_self, [&](auto &s) {
                    s.id= expire_table.available_primary_key();
                    s.trxid= it->trxid;
                    s.timestamp = now();
                    s.producer = it->producer;
                    s.category = it->category;
                });

                it = idx.erase(it);

                ++count;

            }else{
                ++it;
            }
        } 
    }

    void extract::rollbacktrx(){
        records trans_table(_self,_self);
        successtrxs success_table(_self,_self);

        for ( auto it = trans_table.cbegin(); it != trans_table.cend(); ++it) { 
            auto trxid = it->trxid;

            //get transaction details and rollback
            uint64_t subtime = current_time() - it->countdown_time;

            if ((it->category == "DETH" && subtime > 30*60) || (it->category == "DEOS" && subtime > 5*60)){
                        
                transrecords transrecord_table(N(datxos.recharge),N(datxos.recharge));
                auto idx0 = transrecord_table.template get_index<N(fixed_key)>();
                auto itr0 = idx0.find( get_fixed_key(trxid) );

                datxos_assert(itr0 != idx0.end(), "trxid not exists in token transfer");

                // pay the cost
                action(permission_level{ _self, N(active) },
                    N(eosio.token), N(transfer),
                    std::make_tuple(N(datxos.recharge),itr0 -> account, itr0 -> quantity, std::string(""))
                ).send();

                success_table.emplace(_self, [&](auto &s) {
                    s.id= success_table.available_primary_key();
                    s.trxid= trxid;
                    s.producer = it -> producer;
                    s.timestamp = now();
                    s.category = it->category;
                });
                trans_table.erase(*it);

            }
        }
    }

} // namespace Datxos

DATXOS_ABI( datxos::extract, (recordtrx)(setverifiers)(setdoing)(setsuccess)(updateexpire))