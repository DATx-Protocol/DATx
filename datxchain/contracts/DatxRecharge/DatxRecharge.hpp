#pragma once
#include <DatxosLib/DatxosLib.hpp>
#include <DatxosLib/asset.hpp>
#include <string>
namespace datxos
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

    /// @abi action
    void recorduser(account_name datxaddress,
                    string address,
                    account_name bpname);


    /// @abi action
    void charge(account_name bpname,
                string hash,
                string from,
                string to,
                int64_t blocknum,
                string quantity,
                string category,
                string memo);
    public:
     /// @abi table
     struct user
     {
       uint64_t id; //primary key
       checksum256 hash;
        account_name datxaddress; 
        string address;
        
        uint64_t primary_key() const { return id; }
        key256 by_hash() const {return get_hash(hash);}
        
        DATXLIB_SERIALIZE(user, (id)(hash)(datxaddress)(address));
     };

    typedef multi_index<N(user), user,
                        indexed_by<N(hash), const_mem_fun<user, key256, &user::by_hash>>> users;

     /// @abi table
    struct record
    {
        uint64_t id; //primary key
        string trxid;
        string from;
        string to;
        int64_t blocknum;
        string quantity;
        string category;
        string memo;
        checksum256 data;
        uint32_t start_time;
        vector<account_name> producers;
        uint64_t primary_key() const { return id; }
        key256 by_data() const {return get_hash(data);}
        time by_start_time() const {return start_time;}

        DATXLIB_SERIALIZE(record, (id)(trxid)(from)(to)(blocknum)(quantity)(category)(memo)(data)(start_time)(producers));
    };

    typedef multi_index<N(record), record,
                        indexed_by<N(data), const_mem_fun<record, key256, &record::by_data>>> records;

    /// @abi table
        struct countrecord
        {
            uint64_t id; //primary key
            checksum256 data;
            int8_t count;
            uint64_t primary_key() const { return id; }
            key256 by_data() const {return get_hash(data);}
            DATXLIB_SERIALIZE(countrecord, (id)(data)(count));
        };

        typedef multi_index<N(countrecord),countrecord,
                            indexed_by<N(data), const_mem_fun<countrecord, key256, &countrecord::by_data>>> countrecords;



    ///@abi table
    struct expiration
    {
        uint64_t id; //primary key
        string trxid;
        string from;
        string to;
        int64_t blocknum;
        string quantity;
        string category;
        string memo;
        checksum256 data;
        uint64_t primary_key() const { return id; }
        key256 by_data() const {return get_hash(data);}
        DATXLIB_SERIALIZE(expiration, (id)(trxid)(from)(to)(blocknum)(quantity)(category)(memo));
    };

    typedef multi_index<N(expiration), expiration,
                        indexed_by<N(data), const_mem_fun<expiration, key256, &expiration::by_data>>> expirations;


    ///@abi table
     struct perfection
    {
        uint64_t id; //primary key
        string trxid;
        string from;
        string to;
        int64_t blocknum;
        string quantity;
        string category;
        string memo;
        checksum256 data;
        std::vector<account_name> producers;
        uint64_t primary_key() const { return id; }
        key256 by_data() const {return get_hash(data);}
        DATXLIB_SERIALIZE(perfection, (id)(trxid)(from)(to)(blocknum)(quantity)(category)(memo)(data)(producers));
    };

    typedef multi_index<N(perfection), perfection,
                        indexed_by<N(data), const_mem_fun<perfection, key256, &perfection::by_data>>> perfections;


};



}; // namespace Datxos
