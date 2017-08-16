package com.sim.lbox;

import java.util.List;

/**
 * Created by ariel on 16.08.2017.
 */
public class LExpression {
    public String[] input;
    public String calculation;

    public LExpression(String[] arguments, String process) {
        input = arguments;
        calculation = process;
    }
}
