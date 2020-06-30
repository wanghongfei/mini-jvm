package com.fh;
import cn.minijvm.io.Printer;

public class RecursionTest {
    public static void main(String[] args) {
        // 从1打印到100
        foo(1);
    }

    public static void foo(int i) {
        if (i > 100) {
            return;
        }

        Printer.print(i);
        i++;
        foo(i);
    }
}