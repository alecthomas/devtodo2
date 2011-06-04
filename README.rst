DevTodo2
========
DevTodo2 is a command-line task management utility. Tasks are hierarchically
organised, have priorities, and track creation and completion time.

Task lists are stored in the current directory as the file ``.todo2``.

For much more complete information please refer to the man page (todo2(1)).

Installing
----------
DevTodo2 is written in `Go <http://golang.org>`_. To install, you will
need a recent version.

Once you have this, you should be able to simply type::

  $ goinstall github.com/alecthomas/devtodo2

If this fails, try::

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
Yes, this is version 2. `Version 1 <http://swapoff.org/DevTodo>`_ was written in
C++ in 2004, and has been due for a rewrite for a very long time.

Differences between version 1 and 2
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

New features:

- Task lists are now stored as JSON.
- Everything is a *lot* faster.
- Much less code.

Not currently supported:

- Readline-based editing of task text and priority.
- Filters. Completed tasks are hidden by default, but may be displayed with -A.
- Sorting. Tasks are always shown in priority order.
- Linked files.
- Backups.
- ~/.todorc configuration file.
- Colour customisation.
- Custom task formatting.

How do I import my version 1 task lists?
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
DevTodo2 will load your old ``.todo`` files automatically. If you modify the
task list with DevTodo2 it will transparently migrate the file to ``.todo2`` in
the new format.

You can specify the version 1 filename to load with
``--legacy-file=<filename>``, and the version 2 filename with
``--file=<filename>``.
