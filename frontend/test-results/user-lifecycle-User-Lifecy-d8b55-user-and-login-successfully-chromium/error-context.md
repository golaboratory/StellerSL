# Page snapshot

```yaml
- generic [ref=e5]:
  - generic [ref=e6]:
    - generic [ref=e8]: StellerSL
    - generic [ref=e10]: Create an account
  - generic [ref=e12]:
    - generic [ref=e13]:
      - generic [ref=e14]: Full Name
      - textbox "Full Name" [ref=e15]:
        - /placeholder: John Doe
        - text: Test User
    - generic [ref=e16]:
      - generic [ref=e17]: Email
      - textbox "Email" [ref=e18]:
        - /placeholder: email@example.com
        - text: test-1772463742231@example.com
    - generic [ref=e19]:
      - generic [ref=e20]: Password
      - generic [ref=e21]:
        - textbox "Min 8 characters" [ref=e22]: password123
        - img [ref=e23]
        - generic: Enter a password
    - alert [ref=e25]:
      - generic [ref=e28]: Failed to register. Please try again.
    - button "Sign Up" [ref=e29] [cursor=pointer]:
      - generic [ref=e30]: Sign Up
    - generic [ref=e31]:
      - text: Already have an account?
      - link "Login" [ref=e32] [cursor=pointer]:
        - /url: /login
```