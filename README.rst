DevTodo2
========
DevTodo2 is a command-line task management utility. Tasks are hierarchically
organised, have priorities, and track creation and completion time.

Task lists are stored in separate files, allowing for per-project sets of tasks.

For much more complete information please refer to the man page (todo2(1)).

Installing
----------
DevTodo2 is written in `Go <http://golang.org>`_.

To install, you will need a recent version of Go. Once you have this, you should
be able to simply type::

  $ goinstall github.com/alecthomas/devtodo2

If this fails, you can try with::

  $ git clone git://github.com/alecthomas/devtodo2.git
  $ cd devtodo2
  $ gomake
  $ gomake install

Examples
--------
Add a new task::

  $ todo2 -a Shopping list

Add a new sub-task below task 1::

  $ todo2 -ag 1 Buy soap

List outstanding tasks::

  $ todo2

List *all* tasks::

  $ todo2 -A

DevTodo1?
---------
Yes, this is version 2. Version 1 was written in C++ in 2004, and has been due
for a rewrite for a very long time.

DevTodo2 will migrate your old ``.todo`` files automatically.

What isn't supported in version 2 (yet) that was in version 1:

- Readline-based editing of task text and priority is not supported.
- Filters are not yet supported.
- Sorting is not supported. Tasks are always shown in priority order.
- Linked files are not supported.
- Backups are not supported.
- ~/.todorc configuration file is not supported.
- Colour customisation is not supported.
- Custom task formatting is not supported.

An Aside on Go
--------------
As an aside, writing this in Go has mostly been a real pleasure. I normally
code in Python when given a choice, so Go is a nice compromise between the
expressiveness of Python and the performance of C++.

It also has the nice property of generating completely self-contained binaries,
which eradicates one of the hassles of distributing Python applications.

On the downside, some things are much more verbose in Go than they are in
Python. This is mostly due to a lack of exceptions. For example, the following
code in Python will expand a task range (eg. 1.2.3-10)::

	def expand_range(index):
	  start_index, end = index.split('-')
	  start_index, start = start_index.rsplit('.', 1)
	  for i in range(int(start), int(end) + 1):
	    yield '%s.%s' % (start_index, str(i))
  
This is *vastly* more complicated in Go.
