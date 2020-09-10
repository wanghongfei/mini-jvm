
package com.fh;

import cn.minijvm.io.Printer;

public class ObjectTest {
    public static void main(String[] args) throws Exception {
        Object o = new Object();
        Printer.print(o.hashCode());
        Printer.print(o.hashCode());

        MyObject o2 = new MyObject();
        Printer.print(o2.hashCode());
        Object o3 = o2.clone();
        Printer.print(o3.hashCode());
    }

    public static class MyObject implements Cloneable {
        public Object clone() throws CloneNotSupportedException {
            return super.clone();
        }
    }
}