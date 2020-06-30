package com.fh;
import cn.minijvm.io.Printer;

public class IfTest {
    public static void main(String[] args) {
        int sum = 0;

        if (sum > 0) {
            sum = 100;
        }
        if (sum >= 0) {
            sum -= 200;
        }
        if (sum < 0) {
            sum -= 100;
        }
        if (sum <= 0) {
            sum -= 1;
        }
        if (sum == 0) {
            sum += 2000;
        }

        Printer.print(sum); // -301

        if (sum > 10) {

        }
        if (sum >= 10) {

        }
        if (sum < 10) {

        }
        if (sum <= 10) {

        }
        if (sum == 10) {

        }
        Printer.print(sum);
    }

}