#include "vendor/unity.h"
#include <time.h>

#include "vec.h"

extern data_t dotproduct(vec_ptr, vec_ptr);
extern data_t dotproduct1(vec_ptr, vec_ptr);
extern data_t dotproduct2(vec_ptr, vec_ptr);
extern data_t dotproduct3(vec_ptr, vec_ptr);
extern data_t dotproduct4(vec_ptr, vec_ptr);

void setUp(void) {
}

void tearDown(void) {
}

void test_empty(void) {
  vec_ptr u = new_vec(0);
  vec_ptr v = new_vec(0);

  TEST_ASSERT_EQUAL(0, dotproduct(u, v));

  free_vec(u);
  free_vec(v);
}

void test_basic(void) {
  vec_ptr u = new_vec(3);
  vec_ptr v = new_vec(3);

  set_vec_element(u, 0, 1);
  set_vec_element(u, 1, 2);
  set_vec_element(u, 2, 3);
  set_vec_element(v, 0, 4);
  set_vec_element(v, 1, 5);
  set_vec_element(v, 2, 6);

  TEST_ASSERT_EQUAL(32, dotproduct(u, v));

  free_vec(u);
  free_vec(v);
}

void test_longer(void) {
  long n = 1000000;
  vec_ptr u = new_vec(n);
  vec_ptr v = new_vec(n);

  for (long i = 0; i < n; i++) {
    set_vec_element(u, i, i + 1);
    set_vec_element(v, i, i + 1);
  }

  long expected = (2 * n * n * n + 3 * n * n + n) / 6;
  TEST_ASSERT_EQUAL(expected, dotproduct(u, v));

  free_vec(u);
  free_vec(v);
}

void profile(char *fn_name, data_t (*fn)(vec_ptr u, vec_ptr v)) {
  clock_t clock_start, clock_end;
  double clocks_elapsed, time_elapsed, total_ticks = 0.0;

  int runs = 100;
  long n = 1000000;
  for (int i = 0; i < runs; i++) {
    vec_ptr u = new_vec(n);
    vec_ptr v = new_vec(n);

    for (long i = 0; i < n; i++) {
      set_vec_element(u, i, i + 1);
      set_vec_element(v, i, i + 1);
    }

    clock_start = clock();
    (*fn)(u, v);
    clock_end = clock();

    clocks_elapsed = clock_end - clock_start;
    total_ticks += clocks_elapsed;

    free_vec(u);
    free_vec(v);
  }

  time_elapsed = total_ticks / CLOCKS_PER_SEC;
  printf("%s\t%d\t%.0f\t%.3f\t%.4f\t%ld\n",
      fn_name,
      runs,
      total_ticks,
      time_elapsed,
      time_elapsed/runs,
      n);
}

int main(void) {
    UNITY_BEGIN();

    RUN_TEST(test_empty);
    RUN_TEST(test_basic);
    RUN_TEST(test_longer);

    printf("\nFunction\tRuns\tTicks\tTime(s)\tAvg(s)\tVec Len\n");
    profile("dotproduct1", dotproduct1);
    profile("dotproduct2", dotproduct2);
    profile("dotproduct3", dotproduct3);
    profile("dotproduct4", dotproduct4);
    profile("dotproduct", dotproduct);

    return UNITY_END();
}
