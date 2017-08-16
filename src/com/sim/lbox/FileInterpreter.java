package com.sim.lbox;

import javafx.util.Pair;

import java.io.IOException;
import java.net.URI;
import java.net.URISyntaxException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.Scanner;

/**
 * Created by ariel on 16.08.2017.
 */
public class FileInterpreter {
    public void interpreteAllLine(String arg,Main m) {
        List<String> lines = new ArrayList<String>();
        try {
            lines = Files.readAllLines(Paths.get(new URI(arg)));
        } catch (Exception e) {
            System.out.println(e.getMessage());
            new Scanner(System.in).next();
            System.exit(-1);
        }

        for (String line : lines) {
            Pair<String,Boolean> result = m.interpretLine(line);

            if (result.getValue()){
                System.out.println(result.getKey());
            }
        }
    }
}
