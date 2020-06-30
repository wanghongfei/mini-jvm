
package com.fh;

import cn.minijvm.io.Printer;

public class ArrayTest {
    public static void main(String[] args) {
        int[] nums = new int[]{
            1, 2, 3, 4, 5
        };

        Printer.printInt2(nums[0], nums[1]);
        Printer.printInt2(nums[2], nums[3]);

        char[] chars = new char[] {
                '你', '好', '吗'
        };

        Printer.printChar(chars[1]);
        Printer.printChar(chars[2]);
    }
}