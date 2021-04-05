#include "vec.h"

// Unroll the loop
data_t dotproduct(vec_ptr u, vec_ptr v) {
   data_t sum = 0;
   long i, length = vec_length(u); // we can assume both vectors are same length

   data_t *u_data = get_vec_start(u);
   data_t *v_data = get_vec_start(v);

   for (i = 0; i < length; i+=3) {
      sum += u_data[i] * v_data[i];
      sum += u_data[i+1] * v_data[i+1];
      sum += u_data[i+2] * v_data[i+2];
   }

   for (; i < length; i++) {
      sum += u_data[i] * v_data[i];
   }
   return sum;
}

// Unroll the loop
data_t dotproduct4(vec_ptr u, vec_ptr v) {
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

// Skip bounds checking and index directly into array
data_t dotproduct3(vec_ptr u, vec_ptr v) {
   data_t sum = 0;
   long length = vec_length(u); // we can assume both vectors are same length

   data_t *u_data = get_vec_start(u);
   data_t *v_data = get_vec_start(v);

   for (long i = 0; i < length; i++) {
      sum += u_data[i] * v_data[i];
   }
   return sum;
}

// Calculate length outside the loop
data_t dotproduct2(vec_ptr u, vec_ptr v) {
   data_t sum = 0, u_val, v_val;
   long length = vec_length(u); // we can assume both vectors are same length

   for (long i = 0; i < length; i++) { // we can assume both vectors are same length
        get_vec_element(u, i, &u_val);
        get_vec_element(v, i, &v_val);
        sum += u_val * v_val;
   }
   return sum;
}

// Original Function
data_t dotproduct1(vec_ptr u, vec_ptr v) {
   data_t sum = 0, u_val, v_val;

   for (long i = 0; i < vec_length(u); i++) { // we can assume both vectors are same length
        get_vec_element(u, i, &u_val);
        get_vec_element(v, i, &v_val);
        sum += u_val * v_val;
   }
   return sum;
}
