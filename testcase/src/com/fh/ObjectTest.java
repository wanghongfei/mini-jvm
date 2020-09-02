
package com.fh;

import cn.minijvm.io.Printer;

public class ObjectTest {
    public static void main(String[] args) {
        Object o = new Object();
        Printer.print(o.hashCode());
        Printer.print(o.hashCode());
        // Printer.printString(o.toString());
    }
}