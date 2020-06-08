
package com.fh;

public class ArrayTest {
    public static void main(String[] args) {
        int[] nums = new int[]{
            1, 2, 3, 4, 5
        };

        printInt2(nums[0], nums[1]);
        printInt2(nums[2], nums[3]);

        char[] chars = new char[] {
                '你', '好', '吗'
        };

        printChar(chars[1]);
        printChar(chars[2]);
    }
    public static native void printInt2(int num1, int num2);
    public static native void printChar(char ch);
}