// Inspired from this article: https://dev.to/seanchen1991/implementing-an-lru-cache-in-rust-33pp

// use arrayvec::ArrayVec;

#[derive(Default, Copy, Clone)]
pub struct Entry<T: Copy> {
    val: T,
    prev: usize,
    next: usize,
}

// Const Generics ftw!
pub struct LRUCache<T: Copy, const CAP: usize> {
    entries: [Entry<T>; CAP],
    head: usize,
    tail: usize,
    length: usize,
}

impl<T: Default + Copy, const CAP: usize> Default for LRUCache<T, CAP> {
    fn default() -> Self {

        assert!(CAP < usize::max_value(), "Capacity overflow");

        Self {
            entries: [Default::default(); CAP],
            head: 0,
            tail: 0,
            length: 0,
        }
    }
}

struct IterMut<'a, T: Copy, const CAP: usize> {
    cache: &'a mut LRUCache<T, CAP>,
    pos: usize,
    done: bool,
}

impl<'a, T: Copy, const CAP: usize> Iterator for IterMut<'a, T, CAP> {
    type Item = & mut Entry<T>;

    fn next(&mut self) -> Option<Self::Item> {
        if self.done {
            return None;
        }

        let entry = &mut self.cache.entries[self.pos];
        self.pos = entry.next;

        if self.pos == self.cache.tail {
            self.done = true;
        }

        return Some(entry);
    }
}
