public class HelloNative {
    public static void main(String[] args) {
        int sum = 0;
        for (int ix = 1; ix <= 100; ++ix) {
            sum = add(sum, ix);
        }
        print(sum);
    }

    public static int add(int x, int y) {
        return x + y;
    }

    public static native void print(int num);
}
