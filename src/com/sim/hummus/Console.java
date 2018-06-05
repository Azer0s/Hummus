package com.sim.hummus;

import java.util.Scanner;

public class Console {
    public static void write(String s){
        System.out.print(s);
    }

    public static void writeline(String s){
        System.out.println(s);
    }

    public static String read(){
        return new Scanner(System.in).nextLine();
    }
}
