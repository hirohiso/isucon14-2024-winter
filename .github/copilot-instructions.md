

# Copilot Instructions when performing a code review
## Language
When performing a code review, respond in Japanese.

## Prefixes
When performing a code review, use the following prefixes to indicate the type of comment:
```
[must] : must change 
[imo] : my opinion is this, but it's not mandatory (in my opinion) 
[nits] : minor points (nitpick)
[ask] : question  
[fyi] : reference information
```

## Review Perspectives
When performing a code review, follow our internal review checklist.
### A-Design Perspectives
1. is it designed with layers in mind? 2.
2. is the domain logic properly separated? 3. are dependencies properly managed?
Are dependencies properly managed?

### B-Implementation Perspective
1. is the code highly readable? 2.
2. is there any duplicated code? 3. are variable names unambiguous?
Are variable names unambiguous and clear in intent? 4.
4. are there no redundant comments? 5. is exception handling appropriate?
Is exception handling appropriate?

### C-Security Perspectives
1. are input values properly verified? 2.
2. is confidential information properly managed?

### D-Performance Perspective
1. are performance-impacting areas optimized appropriately? 2.
2. are resources being wasted?

### E-UT Perspective
1. are the tests well documented 2. are the test cases representing the specification?
2. do the test cases represent the specification?
3. are the test cases dependent on external processes?

Please evaluate whether these perspectives have been met.

Sample 

```
A-1: Appropriate 
The modified module concentrates only on the domain logic and is designed with layers in mind.

A-2: Inadequate 
The delivery and payment logic is mixed and there is insufficient separation of domain logic.
```

## 🚫 Prohibited
- English output
- Missing or incorrect prefixes
- Skipping evaluation items
- Formatting violations
- Ambiguating answers (e.g., “looks good,” “some,” etc.)

---.

**All of these instructions must be strictly adhered to. ** 
Violations will result in review results being considered invalid.