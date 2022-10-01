#![allow(unused_variables)]
use std::{fmt::Display, usize};

use crate::LRUCache;

pub struct MGLRU<T: Display + Copy, const LRU_N: usize, const LRU_CAP: usize> {
    lrus: [LRUCache<T, LRU_CAP>; LRU_N],
    length: usize,
}

impl<'a, T: Display + Copy + Default, const LRU_N: usize, const LRU_CAP: usize>
    MGLRU<T, LRU_N, LRU_CAP>
{
    pub fn new() -> Self {
        Self {
            lrus: [LRUCache::<T, LRU_CAP>::new(); LRU_N],
            length: 0,
        }
    }

    pub fn len(&self) -> usize {
        return self.length;
    }

    pub fn clear(&mut self) {
        for i in 0..LRU_N {
            self.lrus[i].clear();
            self.length = 0;
        }
    }

    pub fn insert(&mut self, val: T) -> Option<T> {
        let mut deleted_element_opt = Some(val);

        for i in 0..LRU_N {
            deleted_element_opt = match self.lrus[i].insert(deleted_element_opt.unwrap()) {
                Some(entry) => Some(entry.val),
                None => None,
            };

            if deleted_element_opt.is_none() {
                break;
            }
        }

        if self.length < (LRU_N * LRU_CAP) {
            self.length += 1;
        }

        return deleted_element_opt;
    }

    pub fn iter(&'a self) -> Iter<T, LRU_N, LRU_CAP> {
        Iter::new(self)
    }

    fn iter_internal(&'a self) -> IterInternal<T, LRU_N, LRU_CAP> {
        IterInternal::new(self)
    }

    pub fn front(&self) -> Option<&T> {
        if self.length > 0 {
            return None;
        }

        return self.lrus[0].front();
    }

    pub fn front_mut(&mut self) -> Option<&mut T> {
        if self.length > 0 {
            return None;
        }

        return self.lrus[0].front_mut();
    }

    fn remove(&mut self, lru_idx: usize, ele_idx: usize) {
        assert!(self.length != 0);

        self.lrus[lru_idx].remove(ele_idx);
    }

    pub fn touch<F>(&mut self, mut pred: F) -> bool
    where
        F: FnMut(&T) -> bool,
    {
        let mut lru_idx = usize::MAX;
        let mut ele_idx = usize::MAX;
        let mut ele = None;

        self.iter_internal().find(|&result| {
            if pred(result.2) {
                lru_idx = result.0;
                ele_idx = result.1;
                ele = Some(result.2.clone());

                return true;
            }

            false
        });

        if lru_idx != usize::MAX && ele_idx != usize::MAX {
            self.remove(lru_idx, ele_idx);
            self.insert(ele.unwrap());

            return true;
        }

        false
    }

    // For debugging purposes
    // pub fn print_lrus(&self) {
    //     println!("Printing Values in MGLRU");
    //     for i in 0..LRU_N {
    //         println!("LRU {}:", i);
    //         self.lrus[i].print_entries();
    //     }
    // }
}

pub struct Iter<'a, T: Display + Copy + Default, const LRU_N: usize, const LRU_CAP: usize> {
    pos: usize,
    iters: Vec<crate::Iter<'a, T, LRU_CAP>>,
}

impl<'a, T: Display + Copy + Default, const LRU_N: usize, const LRU_CAP: usize>
    Iter<'a, T, LRU_N, LRU_CAP>
{
    fn new(cache: &'a MGLRU<T, LRU_N, LRU_CAP>) -> Iter<'a, T, LRU_N, LRU_CAP> {
        let mut iters = vec![];

        for i in 0..LRU_N {
            iters.push(cache.lrus[i].iter());
        }

        Self { pos: 0, iters }
    }
}

impl<'a, T: Display + Copy + Default, const LRU_N: usize, const LRU_CAP: usize> Iterator
    for Iter<'a, T, LRU_N, LRU_CAP>
{
    type Item = (usize, &'a T);

    fn next(&mut self) -> Option<Self::Item> {
        if self.pos == LRU_N {
            return None;
        }

        let mut next_opt = self.iters[self.pos].next();
        while next_opt.is_none() {
            if self.pos >= self.iters.len() - 1 {
                break;
            }

            self.pos += 1;
            next_opt = self.iters[self.pos].next();
        }

        return next_opt;
    }
}

pub struct IterInternal<'a, T: Display + Copy + Default, const LRU_N: usize, const LRU_CAP: usize> {
    pos: usize,
    iters: Vec<crate::Iter<'a, T, LRU_CAP>>,
}

impl<'a, T: Display + Copy + Default, const LRU_N: usize, const LRU_CAP: usize>
    IterInternal<'a, T, LRU_N, LRU_CAP>
{
    fn new(cache: &'a MGLRU<T, LRU_N, LRU_CAP>) -> IterInternal<'a, T, LRU_N, LRU_CAP> {
        let mut iters = vec![];

        for i in 0..LRU_N {
            iters.push(cache.lrus[i].iter());
        }

        Self { pos: 0, iters }
    }
}

impl<'a, T: Display + Copy + Default, const LRU_N: usize, const LRU_CAP: usize> Iterator
    for IterInternal<'a, T, LRU_N, LRU_CAP>
{
    // Here first usize is LRU number, second usize is index in that LRU
    type Item = (usize, usize, &'a T);

    fn next(&mut self) -> Option<Self::Item> {
        if self.pos == LRU_N {
            return None;
        }

        let mut next_opt = self.iters[self.pos].next();
        while next_opt.is_none() {
            if self.pos >= self.iters.len() - 1 {
                break;
            }

            self.pos += 1;
            next_opt = self.iters[self.pos].next();
        }

        if let Some((ele_idx, ele)) = next_opt {
            return Some((self.pos, ele_idx, ele));
        }

        return None;
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn items<T: Display + Copy + Default, const LRU_N: usize, const LRU_CAP: usize>(
        cache: &MGLRU<T, LRU_N, LRU_CAP>,
    ) -> Vec<T> {
        cache
            .iter()
            .map(|(_, x)| {
                return x.clone();
            })
            .collect()
    }

    #[test]
    fn test_empty() {
        let mut cache = MGLRU::<i32, 2, 1>::new();

        assert_eq!(cache.len(), 0);
        assert_eq!(items(&mut cache), []);
    }

    #[test]
    fn test_basic_insert_order_and_touch() {
        let mut cache = MGLRU::<i32, 2, 1>::new();

        cache.insert(1);

        assert_eq!(cache.len(), 1);
        assert_eq!(items(&cache), [1]);

        cache.insert(2);

        assert_eq!(cache.len(), 2);
        assert_eq!(items(&cache), [2, 1]);

        cache.touch(|x| *x == 1);

        assert_eq!(cache.len(), 2);
        assert_eq!(items(&cache), [1, 2]);
    }

    #[quickcheck]
    fn touch(num: i32) {
        let first = num;
        let second = num + 1;
        let third = num + 2;
        let fourth = num + 3;

        let mut cache = MGLRU::<i32, 2, 2>::new();

        cache.insert(first);
        cache.insert(second);
        cache.insert(third);
        cache.insert(fourth);

        cache.touch(|x| *x == fourth + 1);

        assert_eq!(
            items(&cache),
            [fourth, third, second, first],
            "Nothing is touched."
        );

        cache.touch(|x| *x == second);

        assert_eq!(
            items(&cache),
            [second, fourth, third, first],
            "Touched item is moved to front."
        );
    }

    #[test]
    fn test_clear() {
        let mut cache = MGLRU::<i32, 4, 1>::new();

        cache.insert(1);
        cache.clear();

        assert_eq!(cache.len(), 0);
        assert_eq!(items(&mut cache), [], "All items evicted");

        cache.insert(1);
        cache.insert(2);
        cache.insert(3);
        cache.insert(4);
        assert_eq!(items(&mut cache), [4, 3, 2, 1]);
        cache.clear();
        assert_eq!(items(&mut cache), [], "All items evicted again");
    }
}
