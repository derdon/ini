INI -- a Go library to parse, edit and create ini configuration files
=====================================================================

Installation
------------

::

    go get github.com/derdon/ini

Updating
--------

::

    go get -u github.com/derdon/ini

Usage
-----

See the Wiki at https://github.com/derdon/ini/wiki/ for examples and
explanations on how to use this library.

The API documentation can be found at http://godoc.org/github.com/derdon/ini.

Supported Format
----------------

ini uses a line-based parser. A line is a string that end with a newline
``\n``. If a parsed line is empty, i.e. if it does not even contain
``\n``, the parser reached EOF (end of file) and parsing stops. Lines
consisting only whitespace are ignored. Each non-empty line represents one
*element*. Before and after each element any number of whitespace is
allowed. The supported elements are:

    - comments
    - sections
    - assignments

Comments
~~~~~~~~

A comment begins with a hash sign ``#`` or semicolon ``;``. Before the
introducing comment sign, only whitespace is allowed. That means that the
following line will be parsed to an item with the property ``property``
and the value ``value # this is not a comment``::

    property = value # this is not a comment

Sections
~~~~~~~~

A section begins with an open bracket ``[`` and ends with a closing
bracket ``]``. Between those brackets, there must be at least one
character to name this section. Sections may not be nested!

Assignments
~~~~~~~~~~~

An assignment is internally stored as an item consisting of a *property*
and a *value*. Currently, only the equal sign ``=`` is supported to assign
values to properties. Whitespace before and after the assignment sign is
ignored.

Properties
``````````

A property is a string that starts with a non-whitespace character and
ends with the last non-whitespace character that can be found before the
assignment sign.

Values
``````

A value is a string that starts with the first non-whitespace character
after the assignment sign and ends with the end of the line.

Bugs
----

Quoted values are not supported yet. See TODO.rst in this
folder for more ideas which have not been implemented yet.
