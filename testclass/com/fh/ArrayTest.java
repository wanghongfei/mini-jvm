
package com.fh;

public class ArrayTest {
    public static void main(String[] args) {
        int[] nums = new int[]{
            1, 2, 3, 4, 5
        };

        print(nums[0]);
        print(nums[1]);
    }
    public static native void print(int num);
}