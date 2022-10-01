# LRUST

<tr>
Experiemental LRU and MGLRU ( Muti-Generational LRU ) in Rust

Get aways from LWN articles ( listed further down in the README ) for a general impl of MGLRU:
- MGLRU in kernel is very different and complex due to it being controlled by a variety of factors, small eg. access type file-based page or DMA by CPU, we will try to implement a much simpler version demonstrating just generational use
- There have to be two processes, one which marks what generation this data will be destined for, and other which actually changes the generation ( if it sounds similar to mark and sweep GC, then yes it is. In fact I think the whole idea of MGLRU is inspired from generational GCs )
- By default, new data access will land in younger generation
- A page's generation reflects its "age" — how long it has been since the page was last accessed.
- It's not yet merged as of 5.19 ( source: https://www.phoronix.com/news/MGLRU-v13-More-Benchmarks )

## LWN articles
- https://lwn.net/Articles/851184/
- https://lwn.net/Articles/856931/
- https://lwn.net/Articles/881675/
- https://lwn.net/Articles/894859/

## More info on MGLRU:
- https://github.com/hakavlad/mg-lru-helper

## Notes
- For a bit more insight into decisions and why something is happening, you can browse golang code, as that was the playground for me ( GC langs make it easier to implement such stuff )
- This is a purely experiemental library, if you see improvements which can be made, please do reach out!
- If you try to run `cargo test`, and see some tests are failing, that maybe because of the `quickcheck` crate. For some reason it fails when ran in bulk with other tests. But if you run same tests in isolation ( i.e. one by one ), they will pass.
