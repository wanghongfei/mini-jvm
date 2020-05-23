package com.fh;

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

        print(sum); // -301

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
        print(sum);
    }

    public static native void print(int num);
}