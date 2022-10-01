// Inspired from this article: https://dev.to/seanchen1991/implementing-an-lru-cache-in-rust-33pp

use std::fmt::Display;

mod mglru;

#[cfg(test)]
extern crate quickcheck;
#[cfg(test)]
#[macro_use(quickcheck)]
extern crate quickcheck_macros;

#[derive(Default, Copy, Clone)]
pub struct Entry<T: Copy + Display> {
    pub val: T,
    prev: usize,
    next: usize,
}

impl<T: Copy + Display> Display for Entry<T> {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "Val: {}, prev: {}, next: {}",
            self.val, self.prev, self.next
        )
    }
}

// Const Generics ftw!
#[derive(Clone, Copy)]
pub struct LRUCache<T: Copy + Display, const CAP: usize> {
    entries: [Entry<T>; CAP],
    head: usize,
    tail: usize,
    length: usize,
}

impl<'a, T: Copy + Default + Display, const CAP: usize> LRUCache<T, CAP> {
    fn new() -> Self {
        assert!(CAP < usize::max_value(), "Capacity overflow");

        Self {
            entries: [Default::default(); CAP],
            head: 0,
            tail: 0,
            length: 0,
        }
    }

    pub fn len(&self) -> usize {
        self.length
    }

    pub fn is_empty(&self) -> bool {
        self.length == 0
    }

    pub fn iter_mut(&'a mut self) -> IterMut<'a, T, CAP> {
        IterMut::new(self)
    }

    pub fn iter(&'a self) -> Iter<'a, T, CAP> {
        Iter::new(self)
    }

    pub fn clear(&mut self) {
        self.entries = [Default::default(); CAP];

        self.head = 0;
        self.tail = 0;
        self.length = 0;
    }

    pub fn insert(&mut self, val: T) -> Option<Entry<T>> {
        let entry = Entry {
            val,
            prev: 0,
            next: 0,
        };

        let mut deleted_entry = None;

        let new_head = if self.length == CAP {
            let last_index = self.pop_back();

            deleted_entry = Some(self.entries[last_index]);

            // overwrite the oldest entry with the new entry
            self.entries[last_index] = entry;

            // return the index of the newly-overwritten entry
            last_index
        } else {
            // Here we are only checking lower bound, as
            // upper bound is checked by CAP in if

            if self.length > 0 {
                self.entries[self.length] = entry;
            } else {
                self.entries[0] = entry;
            }

            self.length += 1;

            self.len() - 1
        };

        self.push_front(new_head);

        return deleted_entry;
    }

    /// Returns a reference to the element stored at
    /// the head of the list
    pub fn front(&self) -> Option<&T> {
        Some(&self.entries.first()?.val)
    }

    /// Returns a mutable reference to the element stored
    /// at the head of the list
    pub fn front_mut(&mut self) -> Option<&mut T> {
        Some(&mut self.entries.first_mut()?.val)
    }

    /// Takes an entry that has been added to the linked
    /// list and moves the head to the entry’s position
    fn push_front(&mut self, index: usize) {
        if self.length == 0 {
            return;
        }

        if self.length == 1 {
            self.tail = index;
        } else {
            self.entries[index].next = self.head;
            self.entries[self.head].prev = index;
            // self.entries[self.entries[index].prev].next =
            self.head = index;
        }
    }

    /// Remove the last entry from the list and returns
    /// the index of the removed entry. Note that this
    /// only unlinks the entry from the list, it doesn’t
    /// remove it from the array.
    fn pop_back(&mut self) -> usize {
        let old_tail = self.tail;
        self.tail = self.entries[old_tail].prev;
        old_tail
    }

    pub fn remove(&mut self, index: usize) {
        assert!(self.length != 0);

        let prev = self.entries[index].prev;
        let next = self.entries[index].next;

        if index == self.head {
            self.head = next;
        } else {
            self.entries[prev].next = next;
        }

        if index == self.tail {
            self.tail = prev;
        } else {
            self.entries[next].prev = prev;
        }
    }

    /// Touch a given entry at the given index, putting it
    /// first in the list.
    fn touch_index(&mut self, index: usize) {
        if self.head != index {
            self.remove(index);

            self.length += 1;
            self.push_front(index);
        }
    }

    pub fn touch<F>(&mut self, mut pred: F) -> bool
    where
        F: FnMut(&T) -> bool,
    {
        match self.iter_mut().find(|&(_, ref x)| pred(x)) {
            Some((i, _)) => {
                self.touch_index(i);
                true
            }
            None => false,
        }
    }

    pub fn lookup<F, R>(&mut self, mut pred: F) -> Option<R>
    where
        F: FnMut(&mut T) -> Option<R>,
    {
        let mut result = None;

        for (i, entry) in self.iter_mut() {
            if let Some(r) = pred(entry) {
                result = Some((i, r));
                break;
            }
        }

        match result {
            None => None,
            Some((i, r)) => {
                self.touch_index(i);
                Some(r)
            }
        }
    }

    // For debugging purposes
    pub fn print_entries(&self) {
        println!("LRU Entries: ");

        for entry in self.entries {
            println!("{}", entry);
        }
    }
}

pub struct IterMut<'a, T: Copy + Display, const CAP: usize> {
    cache: &'a mut LRUCache<T, CAP>,
    pos: usize,
    done: bool,
}

impl<'a, T: Copy + Display, const CAP: usize> IterMut<'a, T, CAP> {
    fn new(cache: &'a mut LRUCache<T, CAP>) -> Self {
        let cache_len = cache.length;
        let cache_head = cache.head;

        Self {
            cache,
            pos: cache_head,
            done: cache_len == 0,
        }
    }
}

impl<'a, T: Copy + Display, const CAP: usize> Iterator for IterMut<'a, T, CAP> {
    type Item = (usize, &'a mut T);

    fn next(&mut self) -> Option<Self::Item> {
        if self.done {
            return None;
        }

        let entry = unsafe { &mut *(&mut self.cache.entries[self.pos] as *mut Entry<T>) };
        let index = self.pos;

        if self.pos == self.cache.tail {
            self.done = true;
        }

        self.pos = entry.next;

        Some((index, &mut entry.val))
    }
}

pub struct Iter<'a, T: Copy + Display, const CAP: usize> {
    cache: &'a LRUCache<T, CAP>,
    pos: usize,
    done: bool,
}

impl<'a, T: Copy + Display, const CAP: usize> Iter<'a, T, CAP> {
    fn new(cache: &'a LRUCache<T, CAP>) -> Self {
        let cache_len = cache.length;
        let cache_head = cache.head;

        Self {
            cache,
            pos: cache_head,
            done: cache_len == 0,
        }
    }
}

impl<'a, T: Copy + Display, const CAP: usize> Iterator for Iter<'a, T, CAP> {
    type Item = (usize, &'a T);

    fn next(&mut self) -> Option<Self::Item> {
        if self.done {
            return None;
        }

        let entry = &self.cache.entries[self.pos];
        // let entry = unsafe { &mut *(&mut self.cache.entries[self.pos] as *mut Entry<T>) };
        let index = self.pos;

        if self.pos == self.cache.tail {
            self.done = true;
        }

        self.pos = entry.next;

        Some((index, &entry.val))
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn items<T: Copy + Default + Display, const CAP: usize>(
        cache: &mut LRUCache<T, CAP>,
    ) -> Vec<T> {
        cache
            .iter_mut()
            .map(|(_, x)| {
                return x.clone();
            })
            .collect()
    }

    // So if you try to run test file, tests with quickcheck
    // dependency will fail, but if you try to run them
    // as an individual test, they will pass

    #[test]
    fn test_empty() {
        let mut cache = LRUCache::<i32, 1>::new();
        assert_eq!(cache.len(), 0);
        assert_eq!(items(&mut cache), []);
    }

    #[test]
    fn test_insert() {
        let mut cache = LRUCache::<i32, 4>::new();

        cache.insert(1);
        assert_eq!(cache.len(), 1);

        cache.insert(2);
        assert_eq!(cache.len(), 2);

        cache.insert(3);
        assert_eq!(cache.len(), 3);

        cache.insert(4);
        assert_eq!(cache.len(), 4);

        assert_eq!(
            items(&mut cache),
            [4, 3, 2, 1],
            "Ordered from most- to least-recent"
        );

        cache.insert(5);
        assert_eq!(cache.len(), 4);
        assert_eq!(
            items(&mut cache),
            [5, 4, 3, 2],
            "Least-recently-used item evicted"
        );

        cache.insert(6);
        cache.insert(7);
        cache.insert(8);
        cache.insert(9);

        assert_eq!(cache.len(), 4);
        assert_eq!(
            items(&mut cache),
            [9, 8, 7, 6],
            "Least-recently-used item evicted"
        );
    }

    #[test]
    fn test_lookup() {
        let mut cache = LRUCache::<i32, 4>::new();

        cache.insert(1);
        cache.insert(2);
        cache.insert(3);
        cache.insert(4);

        let result = cache.lookup(|x| if *x == 5 { Some(()) } else { None });
        assert_eq!(result, None, "Cache miss.");
        assert_eq!(items(&mut cache), [4, 3, 2, 1], "Order not changed.");

        // Cache hit
        let result = cache.lookup(|x| if *x == 3 { Some(*x * 2) } else { None });
        assert_eq!(result, Some(6), "Cache hit.");
        assert_eq!(
            items(&mut cache),
            [3, 4, 2, 1],
            "Matching item moved to front."
        );
    }

    #[test]
    fn test_clear() {
        let mut cache = LRUCache::<i32, 4>::new();

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

    #[quickcheck]
    fn touch(num: i32) {
        let first = num;
        let second = num + 1;
        let third = num + 2;
        let fourth = num + 3;

        let mut cache = LRUCache::<i32, 4>::new();

        cache.insert(first);
        cache.insert(second);
        cache.insert(third);
        cache.insert(fourth);

        cache.touch(|x| *x == fourth + 1);

        assert_eq!(
            items(&mut cache),
            [fourth, third, second, first],
            "Nothing is touched."
        );

        cache.touch(|x| *x == second);

        assert_eq!(
            items(&mut cache),
            [second, fourth, third, first],
            "Touched item is moved to front."
        );
    }

    #[quickcheck]
    fn front(num: i32) {
        let first = num;
        let second = num + 1;

        let mut cache = LRUCache::<i32, 4>::new();

        assert_eq!(cache.front(), None, "Nothing is in the front.");

        cache.insert(first);
        cache.insert(second);

        assert_eq!(
            cache.front(),
            Some(&second),
            "The last inserted item should be in the front."
        );

        cache.touch(|x| *x == first);

        assert_eq!(
            cache.front(),
            Some(&first),
            "Touched item should be in the front."
        );
    }
}
