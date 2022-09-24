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

impl<'a, T: Copy + Default, const CAP: usize> LRUCache<T, CAP> {
    pub fn len(&self) -> usize {
        self.length
    }

    pub fn is_empty(&self) -> bool {
        self.length == 0
    }

    pub fn iter_mut(&'a mut self) -> IterMut<'a, T, CAP> {
        IterMut::new(self)
    }

    pub fn clear(&mut self) {
        self.entries = [Default::default(); CAP];

        self.head = 0;
        self.tail = 0;
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

    fn remove(&mut self, index: usize) {
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

pub struct IterMut<'a, T: Copy, const CAP: usize> {
    cache: &'a mut LRUCache<T, CAP>,
    pos: usize,
    done: bool,
}

impl<'a, T: Copy, const CAP: usize> IterMut<'a, T, CAP> {
    fn new(cache: &'a mut LRUCache<T, CAP>) -> Self {
        let cache_len = cache.length;

        Self {
            cache,
            pos: 0,
            done: cache_len == 0,
        }
    }
}

impl<'a, T: Copy, const CAP: usize> Iterator for IterMut<'a, T, CAP> {
    type Item = (usize, &'a mut T);

    fn next(&mut self) -> Option<Self::Item> {
        if self.done {
            return None;
        }

        let entry = unsafe { &mut *(&mut self.cache.entries[self.pos] as *mut Entry<T>) };
        let index = self.pos;

        if self.cache.tail == self.pos {
            self.done = true;
        }

        self.pos = entry.next;

        Some((index, &mut entry.val))
    }
}
