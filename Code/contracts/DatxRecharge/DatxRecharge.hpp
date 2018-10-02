#pragma once
#include <DatxioLib/DatxioLib.hpp>
#include <DatxioLib/asset.hpp>

#include <string>

namespace datxio
{
using std::string;
class recharge : public contract
{
  public:
    recharge(account_name self) : contract(self) {}

    /// @abi action
    void charge(transaction_id_type hash,
                string from,
                string to,
                int64_t blocknum,
                string quantity,
                string category,
                string memo);

    /// @abi action
    void transtoken(transaction_id_type hash,
            account_name from,
            account_name to,
            asset quantity,
            string memo);

  public:
     /// @abi table
    struct record
    {
        uint64_t trxid; //primary key
        transaction_id_type hash;
        string from;
        string to;
        int64_t blocknum;
        string quantity;
        string category;
        string memo;

        uint64_t primary_key() const { return trxid; }
        key256 by_hash() const {return get_hash(hash);}
        static key256 get_hash(const checksum256& hash) {
            const uint64_t *p64 = reinterpret_cast<const uint64_t *>(&hash);
            return key256::make_from_word_sequence<uint64_t>(p64[0], p64[1], p64[2], p64[3]);
         }
        DATXLIB_SERIALIZE(record, (trxid)(hash)(from)(to)(quantity)(blocknum)(category)(memo))
    };
    /// @abi table
    struct transaction
    {
        uint64_t trxid; //primary key
        transaction_id_type hash;
        account_name from;
        account_name to;
        asset quantity;
        string memo;

        uint64_t primary_key() const { return trxid; }
        key256 by_hash() const {return get_hash(hash);}
        static key256 get_hash(const checksum256& hash) {
            const uint64_t *p64 = reinterpret_cast<const uint64_t *>(&hash);
            return key256::make_from_word_sequence<uint64_t>(p64[0], p64[1], p64[2], p64[3]);
         }
        DATXLIB_SERIALIZE(transaction, (trxid)(hash)(from)(to)(quantity)(memo))
    };
    typedef multi_index<N(record), record,indexed_by<N(hash), const_mem_fun<record, key256, &record::by_hash>>> records;
    typedef multi_index<N(transaction), transaction,indexed_by<N(hash), const_mem_fun<transaction, key256, &transaction::by_hash>>> transactions;


};

} // namespace Datxio
