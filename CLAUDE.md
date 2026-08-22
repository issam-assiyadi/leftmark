# Agent instructions

## Writing comments

Default to no comment. Identifiers and structure should carry the meaning; add a comment only when it captures something the code genuinely can't say for itself.

Before adding a comment, check it against all of these:

1. **Not redundant.** If a reader would already understand the line from the code alone, skip the comment. A comment earns its place by explaining *why*, never by restating *what*.
2. **Makes sense out of context, later.** A comment has to stand on its own for someone with no memory of the current task or conversation. Don't write:
   - the story of a bug and how it was found or fixed
   - a comparison to a previous or rejected approach ("unlike the old version, which...")
   - anything shaped like a changelog entry or a narrated decision

   State the current constraint or invariant as a plain fact about the system, not as a narrative about how you got here. "gocui's Clear() resets the view's origin, so it's restored explicitly below" is fine. "We found a bug where scrolling broke because of Clear(), so we fixed it by..." is not.
3. **If you can't tell whether a comment clears the first two checks, ask the user** instead of guessing. Comments are cheap to add and expensive for someone else to clean up later.
4. **Run the `unslop` skill on any comment before committing it**, to strip AI writing tells. If it isn't available, ask the user to install it: `npx skills add https://github.com/cursor/plugins --skill unslop`.
5. **Keep comments in sync with the code.** When a change makes an existing comment inaccurate or its strategy obsolete, update or delete that comment as part of the same change. Don't leave a comment describing an approach the code no longer uses.
