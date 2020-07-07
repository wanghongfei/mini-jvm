package com.fh;
import cn.minijvm.io.Printer;

public class StringTest {
    public static void main(String[] args) {
        String name = "hello, 世界";
        Printer.printString(name);

        String dcs = new String("数字战斗模拟");
        Printer.printString(dcs);
    }
}
