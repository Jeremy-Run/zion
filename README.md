# zion

Hotring implemented by Golang.

Test_Get execution result：
```bash
=== RUN   Test_Get
slot number: 0, link: {key:hello4, value: world4, tag: 1} -> {key:hello7, value: world7, tag: 2} -> HEAD 
slot number: 1, link: {key:hello2, value: world2, tag: 1} -> {key:hello3, value: world3, tag: 2} -> {key:hello5, value: world5, tag: 3} -> HEAD 
slot number: 2, link: {key:hello8, value: world8, tag: 1} -> HEAD 
slot number: 3, link: {key:hello10, value: world10, tag: 1} -> HEAD 
slot number: 4, link: nil 
slot number: 5, link: {key:hello1, value: world1, tag: 1} -> HEAD 
slot number: 6, link: nil 
slot number: 7, link: {key:hello0, value: world0, tag: 1} -> {key:hello6, value: world6, tag: 2} -> {key:hello9, value: world9, tag: 3} -> HEAD 
======================================
key: hello7, value: world7 
======================================
slot number: 0, link: {key:hello7, value: world7, tag: 2} -> {key:hello4, value: world4, tag: 1} -> HEAD 
slot number: 1, link: {key:hello2, value: world2, tag: 1} -> {key:hello3, value: world3, tag: 2} -> {key:hello5, value: world5, tag: 3} -> HEAD 
slot number: 2, link: {key:hello8, value: world8, tag: 1} -> HEAD 
slot number: 3, link: {key:hello10, value: world10, tag: 1} -> HEAD 
slot number: 4, link: nil 
slot number: 5, link: {key:hello1, value: world1, tag: 1} -> HEAD 
slot number: 6, link: nil 
slot number: 7, link: {key:hello0, value: world0, tag: 1} -> {key:hello6, value: world6, tag: 2} -> {key:hello9, value: world9, tag: 3} -> HEAD 
--- PASS: Test_Get (0.00s)
PASS
```

You can see:
```bash
======================================
key: hello7, value: world7 
======================================
```
Above this log is the hashtable originally constructed, and we can see that `key: hello7` was originally ranked in the second node of `slot number: 0`. But under this log, we see that `key:hello7` is in the first node of `slot number: 0`. This is the problem solved by hotring. When a hotkey is accessed multiple times, the probability of its ranking is higher than that of other keys, and the earlier it is placed in the linked list, the faster it can be accessed, which improves the Access efficiency of hashtable.

That's All, Thank you.

