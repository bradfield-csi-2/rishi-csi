#include "vec.h"


data_t dotproduct(vec_ptr u, vec_ptr v) {
   data_t sum = 0;
   long length = vec_length(u); // we can assume both vectors are same length

   data_t *u_data = get_vec_start(u);
   data_t *v_data = get_vec_start(v);

   for (long i = 0; i < length; i++) {
      sum += u_data[i] * v_data[i];
   }
   return sum;
}
