#ifndef _PTHREAD_IMPL_H
#define _PTHREAD_IMPL_H

//include <pthread.h>
//include <signal.h>
#include <errno.h>
#include <limits.h>
#include "libc.h"
//include "syscall.h"
//include "atomic.h"
//include "futex.h"

#define pthread __pthread

struct pthread {
 
   int errno_val;
  
	locale_t locale;

};

#include "pthread_arch.h"

#endif
