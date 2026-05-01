package com.fh;

import cn.minijvm.io.Printer;

public class IfNullTest {
    public static void main(String[] args) {
        Object obj = new Object();

        // ifnull on non-null: should NOT take the branch
        if (obj == null) {
            Printer.print(999);
        } else {
            Printer.print(1);
        }

        obj = null;

        // ifnull on null: should take the branch
        if (obj == null) {
            Printer.print(2);
        } else {
            Printer.print(999);
        }

        // Combined: ifnull + ifnonnull
        String str = "hello";
        if (str == null) {
            Printer.print(999);
        }
        if (str != null) {
            Printer.print(3);
        }

        str = null;
        if (str == null) {
            Printer.print(4);
        }
        if (str != null) {
            Printer.print(999);
        }
    }
}
