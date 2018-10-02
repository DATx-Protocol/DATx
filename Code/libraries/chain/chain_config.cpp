/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */

#include <datxio/chain/chain_config.hpp>
#include <datxio/chain/exceptions.hpp>

namespace datxio { namespace chain {

   void chain_config::validate()const {
      DATX_ASSERT( target_block_net_usage_pct <= config::percent_100, action_validate_exception,
                  "target block net usage percentage cannot exceed 100%" );
      DATX_ASSERT( target_block_net_usage_pct >= config::percent_1/10, action_validate_exception,
                  "target block net usage percentage must be at least 0.1%" );
      DATX_ASSERT( target_block_cpu_usage_pct <= config::percent_100, action_validate_exception,
                  "target block cpu usage percentage cannot exceed 100%" );
      DATX_ASSERT( target_block_cpu_usage_pct >= config::percent_1/10, action_validate_exception,
                  "target block cpu usage percentage must be at least 0.1%" );

      DATX_ASSERT( max_transaction_net_usage < max_block_net_usage, action_validate_exception,
                  "max transaction net usage must be less than max block net usage" );
      DATX_ASSERT( max_transaction_cpu_usage < max_block_cpu_usage, action_validate_exception,
                  "max transaction cpu usage must be less than max block cpu usage" );

      DATX_ASSERT( base_per_transaction_net_usage < max_transaction_net_usage, action_validate_exception,
                  "base net usage per transaction must be less than the max transaction net usage" );
      DATX_ASSERT( (max_transaction_net_usage - base_per_transaction_net_usage) >= config::min_net_usage_delta_between_base_and_max_for_trx,
                  action_validate_exception,
                  "max transaction net usage must be at least ${delta} bytes larger than base net usage per transaction",
                  ("delta", config::min_net_usage_delta_between_base_and_max_for_trx) );
      DATX_ASSERT( context_free_discount_net_usage_den > 0, action_validate_exception,
                  "net usage discount ratio for context free data cannot have a 0 denominator" );
      DATX_ASSERT( context_free_discount_net_usage_num <= context_free_discount_net_usage_den, action_validate_exception,
                  "net usage discount ratio for context free data cannot exceed 1" );

      DATX_ASSERT( min_transaction_cpu_usage <= max_transaction_cpu_usage, action_validate_exception,
                  "min transaction cpu usage cannot exceed max transaction cpu usage" );
      DATX_ASSERT( max_transaction_cpu_usage < (max_block_cpu_usage - min_transaction_cpu_usage), action_validate_exception,
                  "max transaction cpu usage must be at less than the difference between the max block cpu usage and the min transaction cpu usage" );

      DATX_ASSERT( 1 <= max_authority_depth, action_validate_exception,
                  "max authority depth should be at least 1" );
}

} } // namespace datxio::chain