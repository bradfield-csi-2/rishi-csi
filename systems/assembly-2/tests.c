#include "vendor/unity.h"
#include <stdlib.h>

extern int fib(int n);
extern int index(int *matrix, int rows, int cols, int rindex, int cindex);
extern void transpose(int *in, int *out, int rows, int cols);
extern float volume(float radius, float height);

void setUp(void) {
}

void tearDown(void) {
}

void test_fib_0(void) { TEST_ASSERT_EQUAL(0, fib(0)); }
void test_fib_1(void) { TEST_ASSERT_EQUAL(1, fib(1)); }
void test_fib_2(void) { TEST_ASSERT_EQUAL(1, fib(2)); }
void test_fib_3(void) { TEST_ASSERT_EQUAL(2, fib(3)); }
void test_fib_10(void) { TEST_ASSERT_EQUAL(55, fib(10)); }
void test_fib_12(void) { TEST_ASSERT_EQUAL(144, fib(12)); }

void test_index_row(void) {
  int matrix[1][4] = {{1, 2, 3, 4}};
  TEST_ASSERT_EQUAL(3, index((int *)matrix, 1, 4, 0, 2));
}

void test_index_col(void) {
  int matrix[4][1] = {{1}, {2}, {3}, {4}};
  TEST_ASSERT_EQUAL(2, index((int *)matrix, 4, 1, 1, 0));
}

void test_index_rect(void) {
  int matrix[2][3] = {{1, 2, 3}, {4, 5, 6}};
  TEST_ASSERT_EQUAL(6, index((int *)matrix, 2, 3, 1, 2));
}

void test_cone_volume_0_0(void) {
  TEST_ASSERT_FLOAT_WITHIN(0.01, 0.0, volume(0.0, 0.0));
}
void test_cone_volume_1_2(void) {
  TEST_ASSERT_FLOAT_WITHIN(0.01, 2.09, volume(1.0, 2.0));
}
void test_cone_volume_55_55(void) {
  TEST_ASSERT_FLOAT_WITHIN(0.01, 174.23, volume(5.5, 5.5));
}
void test_cone_volume_1234_5678(void) {
  TEST_ASSERT_FLOAT_WITHIN(0.01, 9.05, volume(1.234, 5.678));
}

void test_transpose_wide(void) {
  int in[2][3] = {{1, 2, 3}, {4, 5, 6}};
  int out[3][2];

  transpose((int *)in, (int *)out, 2, 3);

  int expected[3][2] = {{1, 4}, {2, 5}, {3, 6}};
  TEST_ASSERT_EQUAL_INT_ARRAY(expected, out, 6);
}

void test_transpose_tall(void) {
  int in[3][2] = {{1, 2}, {3, 4}, {5, 6}};
  int out[2][3];

  transpose((int *)in, (int *)out, 3, 2);

  int expected[2][3] = {{1, 3, 5}, {2, 4, 6}};
  TEST_ASSERT_EQUAL_INT_ARRAY(expected, out, 6);
}

void test_transpose_square(void) {
  int in[3][3] = {{1, 2, 3}, {4, 5, 6}, {7, 8, 9}};
  int out[3][3];

  transpose((int *)in, (int *)out, 3, 3);

  int expected[3][3] = {{1, 4, 7}, {2, 5, 8}, {3, 6, 9}};
  TEST_ASSERT_EQUAL_INT_ARRAY(expected, out, 9);
}

int main(void) {
    UNITY_BEGIN();

    RUN_TEST(test_fib_0);
    RUN_TEST(test_fib_1);
    RUN_TEST(test_fib_2);
    RUN_TEST(test_fib_3);
    RUN_TEST(test_fib_10);
    RUN_TEST(test_fib_12);

    RUN_TEST(test_index_row);
    RUN_TEST(test_index_col);
    RUN_TEST(test_index_rect);

    RUN_TEST(test_cone_volume_0_0);
    RUN_TEST(test_cone_volume_1_2);
    RUN_TEST(test_cone_volume_55_55);
    RUN_TEST(test_cone_volume_1234_5678);

    RUN_TEST(test_transpose_wide);
    RUN_TEST(test_transpose_tall);
    RUN_TEST(test_transpose_square);

    return UNITY_END();
}
