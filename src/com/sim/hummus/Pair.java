package com.sim.hummus;

public class Pair <T,K> {
    private T key;
    private K value;

    public T getKey() {
        return key;
    }

    public K getValue() {
        return value;
    }

    public Pair(T key, K value){
        this.key = key;
        this.value = value;
    }
}
