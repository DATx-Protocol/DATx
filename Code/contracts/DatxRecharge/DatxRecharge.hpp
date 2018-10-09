#pragma once
#include <DatxioLib/DatxioLib.hpp>
#include <DatxioLib/asset.hpp>

#include <string>

namespace datxio
{
using std::string;
class recharge : public contract
{
  private:
    void expired_trx();
  public:
    recharge(account_name self) : contract(self) {}
    static key256 get_hash(const checksum256& hash) {
            const uint64_t *p64 = reinterpret_cast<const uint64_t *>(&hash);
            return key256::make_from_word_sequence<uint64_t>(p64[0], p64[1], p64[2], p64[3]);
         }
    ///@abi action
    void recorduser(account_name datxaddress,
                    transaction_id_type address,
                    account_name bpname);


    /// @abi action
    void charge(transaction_id_type hash,
                account_name bpname,
                string from,
                string to,
                int64_t blocknum,
                string quantity,
                string category,
                string memo);

     /// @abi table
     struct user
     {
        uint64_t userid; //primary key
        account_name datxaddress; 
        transaction_id_type address;
        uint64_t primary_key() const { return userid; }
        key256 by_address() const {return get_hash(address);}
        account_name get_datxaddress() const {return datxaddress;}
        DATXLIB_SERIALIZE(user, (userid)(datxaddress)(address))
     }

    typedef multi_index<N(user), user,
                        indexed_by<N(address), const_mem_fun<user, key256, &user::by_address>>,
                        indexed_by<N(datxaddress), const_mem_fun<user, account_name, &user::get_datxaddress>>> users;

     /// @abi table
    struct record
    {
        uint64_t id; //primary key
        transaction_id_type trxid;
        string from;
        string to;
        int64_t blocknum;
        string quantity;
        string category;
        string memo;
        checksum256 data;
        time start_time;
        vector<account_name> producers;
        uint64_t primary_key() const { return trxid; }
        key256 by_data() const {return get_hash(data);}
        account_name get_bpname() const { return bpname;}
        time by_start_time() const {return start_time;}

        DATXLIB_SERIALIZE(record, (id)(trxid)(from)(to)(quantity)(blocknum)(category)(memo)(data)(start_time)(producers))
    };

    typedef multi_index<N(record), record,
                        indexed_by<N(data), const_mem_fun<record, key256, &record::by_hash>>
                        ,indexed_by<N(start_time),const_mem_fun<record,time,&record::by_start_time>>> records;

    /// @abi table
        struct countrecord
        {
            uint64_t id; //primary key
            checksum256 data;
            int8_t count;
            uint64_t primary_key() const { return id; }
            key256 by_data() const {return get_hash(data);
            DATXLIB_SERIALIZE(countrecord, (id)(data)(count))
        };

        typedef multi_index<N(countrecord),countrecord,
                            indexed_by<N(data), const_mem_fun<countrecord, key256, &countrecord::by_hash>>> countrecords;



    ///@abi table
    struct expiredrecord
    {
        uint64_t id; //primary key
        transaction_id_type trxid;
        string from;
        string to;
        int64_t blocknum;
        string quantity;
        string category;
        string memo;
        checksum256 data;
        uint64_t primary_key() const { return id; }
        key256 by_data() const {return get_hash(data);}
        account_name get_bpname() const { return bpname;}
        DATXLIB_SERIALIZE(expiredrecord, (id)(trxid)(from)(to)(quantity)(blocknum)(category)(memo))
    };

    typedef multi_index<N(expiredrecord), expiredrecord,
                        indexed_by<N(data), const_mem_fun<expiredrecord, key256, &expiredrecord::by_hash>>> expiredrecords;


    ///@abi table
     struct successrecord
    {
        uint64_t id; //primary key
        transaction_id_type trxid;
        string from;
        string to;
        int64_t blocknum;
        string quantity;
        string category;
        string memo;
        checksum256 data;
        time_point_sec time;
        std::vector<account_name> producers;
        uint64_t primary_key() const { return id; }
        key256 by_data() const {return get_hash(data);}
        DATXLIB_SERIALIZE(successrecord, (id)(trxid)(from)(to)(quantity)(blocknum)(category)(memo)(data)(time)(producers))
    };

    typedef multi_index<N(successrecord), successrecord,
                        indexed_by<N(data), const_mem_fun<successrecord, key256, &successrecord::by_hash>>> successrecords;


};



} // namespace Datxio
