package com.fh;

import cn.minijvm.io.Printer;

public class ForLoopPrintTest {
    public static void main(String[] args) {
        int sum = 0;
        for (int ix = 1; ix <= 100; ++ix) {
            sum = add(sum, ix);
        }

        Printer.print(sum);
    }

    public static int add(int x, int y) {
        return x + y;
    }
}
