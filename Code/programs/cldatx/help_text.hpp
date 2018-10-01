/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once
#include <fc/exception/exception.hpp>

namespace datxio { namespace client { namespace help {
   bool print_recognized_errors(const fc::exception& e, const bool verbose_errors);
   bool print_help_text(const fc::exception& e);
}}}