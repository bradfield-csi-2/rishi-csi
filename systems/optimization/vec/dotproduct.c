#include "vec.h"


data_t dotproduct(vec_ptr u, vec_ptr v) {
   data_t sum = 0;
   long i, length = vec_length(u); // we can assume both vectors are same length

   data_t *u_data = get_vec_start(u);
   data_t *v_data = get_vec_start(v);

   for (i = 0; i < length; i+=2) {
      sum += u_data[i] * v_data[i];
      sum += u_data[i+1] * v_data[i+1];
   }

   for (; i < length; i++) {
      sum += u_data[i] * v_data[i];
   }
   return sum;
}
