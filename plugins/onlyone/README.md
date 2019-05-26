# onlyone

## Name

*onlyone* - Randomly chooses a single RR of each specified type in the response.

## Description

There can be only one. If there are multiple records of the same type in the response,
this plugin will strip out all but one of them. This can be useful, for example, to
replace the *loadbalance* plugin if [Happy Eyeballs sorting](https://tools.ietf.org/html/rfc8305#section-4)
is getting in the way.

## Syntax

~~~
onlyone [ZONES ...] {
  types TYPES
}
~~~

* **ZONES** is the zones to which this plugin applies. If omitted, they default to `.`.
* **TYPES** is a list of DNS types to which this plugin applies. If `types` is omitted, it
  will apply to A and AAAA records only.

## Examples

~~~
onlyone example.com {
  types A SRV AAAA
}
~~~
