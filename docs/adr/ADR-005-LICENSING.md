# ADR-005: Licensing & Copyright

### Submitters
- Jamie Murphy <murpjami@oregonstate.edu>

## Status: Proposed

## Context

This ADR lays out how the Supertext program and any affiliated software programs should be licensed.

## Decision

The Supertext program, along with any affiliated software written in conjunction, will be licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html). This is acceptable
under the terms set out by the client, and the GPLv3 is compatible with the BSD 3-Clause license most Go libraries are written in. GPLv3 is well-tested and retains many legal rights for the software while allowing
others to benefit from our work, while enforcing that all future changes must be made available for all to benefit from.

## Alternatives Considered

One option was to keep the program as closed-source. It was determined that while this would be acceptable, we felt that keeping the program open-source would allow for easier enhancements by future authors
if we decide to stop working on this program after graduation. As well, keeping the program open-source allows for other entities to benefit from our work instead of restricting the work to the client, Portland Opera.

Another option was to not license the program at all. However, this leads the program in a legal grey area, where the legal rights to use or transfer the program for other entities are not clearly defined.
It is important that legal rights be kept and defined, so the program can be used within the necessary contexts.

## Consequences

### Positive

- **Future Modifications Under Same License**: Any future modifications made to this software must be provided under the same license, ensuring future improvements by us or other entities
  can benefit the community as a whole

### Negative

- **License Compatibility**: The GPLv3 is incompatible with GPL-2.0-only, and conditionally compatible with AGPL-3.0 and MPL-2.0. Any libraries or components using these licenses
  is incompatible with this software and must not be used.

### Neutral

- **Commercial and Private Usage**: The GPLv3 allows for commercial and private usage as long as conditions are upheld.

## Changelog

- Rev 1 (2026-10-31) [Jamie Murphy <murpjami@oregonstate.edu>]
  - Initial ADR
