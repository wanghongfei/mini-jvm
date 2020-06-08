
package com.fh;

public class ArrayTest {
    public static void main(String[] args) {
        int[] nums = new int[]{
            1, 2, 3, 4, 5
        };

        printInt(nums[0]);
        printInt(nums[1]);

        char[] chars = new char[] {
                '你', '好', '吗'
        };

        printChar(chars[1]);
        printChar(chars[2]);
    }
    public static native void printInt(int num);
    public static native void printChar(char ch);
}