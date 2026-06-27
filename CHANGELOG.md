## [0.8.0](https://github.com/sceredi/co-type/compare/v0.7.0...v0.8.0) (2026-06-27)

### Features

* generate code snippets from server own code ([8bcea9c](https://github.com/sceredi/co-type/commit/8bcea9cd60675a01a6b52ce5a9ee1edba4d6311b))

### Dependency updates

* **deps:** update node.js to 24.18 ([301993f](https://github.com/sceredi/co-type/commit/301993f2fb937a1944633edb8247e6bd4135623d))

### Bug Fixes

* avoid game window being too wide ([51fe641](https://github.com/sceredi/co-type/commit/51fe64167dafa4118a595e699208ae4d501eac6a))

### Build and continuous integration

* **deps:** update actions/setup-go action to v6.5.0 ([168b832](https://github.com/sceredi/co-type/commit/168b8329749edd2af8889d398f9ab85f5e41c430))

## [0.7.0](https://github.com/sceredi/co-type/compare/v0.6.0...v0.7.0) (2026-06-20)

### Features

* **client:** wire user ready update with server ([2bb8345](https://github.com/sceredi/co-type/commit/2bb8345e0d0024e821f1dd883fa325c35b5d2609))
* **server:** add lobby user ready set ([6017919](https://github.com/sceredi/co-type/commit/601791904c1497d5c5d2fabe0f7e245f00a20189))

### General maintenance

* add lobby service user ready set endpoint ([33da223](https://github.com/sceredi/co-type/commit/33da223136c083f4a77c68aa53f6712302a04633))

## [0.6.0](https://github.com/sceredi/co-type/compare/v0.5.0...v0.6.0) (2026-06-20)

### Features

* **client:** wire user settings update with server ([37e1ba0](https://github.com/sceredi/co-type/commit/37e1ba081e2906e7a17540329e041c7fa264b047))
* **server:** add lobby user update ([a6ad51a](https://github.com/sceredi/co-type/commit/a6ad51afffa227bb36248af1a33e5a1071c27bf7))

### Build and continuous integration

* **deps:** update danysk/action-checkout action to v0.2.30 ([8fc9602](https://github.com/sceredi/co-type/commit/8fc9602c1607d64e3f39cd91be4c4d9b20ab5acf))

### General maintenance

* add lobby service user settings update ([2293d54](https://github.com/sceredi/co-type/commit/2293d543b0da98e6cdae4242532e155f208a3e99))
* remove host from lobby message ([7289111](https://github.com/sceredi/co-type/commit/72891117824f13316cdb3068df2e059531350180))

### Refactoring

* remove host from lobby ([8ed7a0f](https://github.com/sceredi/co-type/commit/8ed7a0f7d0bef9514eb559339c9452a50e324be4))

## [0.5.0](https://github.com/sceredi/co-type/compare/v0.4.0...v0.5.0) (2026-06-19)

### Features

* **broker:** add unregister lobby ([0c7e61c](https://github.com/sceredi/co-type/commit/0c7e61c644ccbfab9635c776dc4980d742a3d14c))
* **client:** add lobby leave ([85e4b2a](https://github.com/sceredi/co-type/commit/85e4b2afb02080708d4259ad57ed9d03e017cc06))
* **server:** add lobby user remove ([3ae9c02](https://github.com/sceredi/co-type/commit/3ae9c022704309ef378fc0550ec95d395450b21c))

### General maintenance

* add broker lobby unregister proto ([335e110](https://github.com/sceredi/co-type/commit/335e1106dc2bfcbd2cdaa688db5661e613b7a55d))
* add lobby service user remove ([997949e](https://github.com/sceredi/co-type/commit/997949ed142897f7529c41446cccc1b43df40e59))

## [0.4.0](https://github.com/sceredi/co-type/compare/v0.3.0...v0.4.0) (2026-06-19)

### Features

* add lobby join ([d76f461](https://github.com/sceredi/co-type/commit/d76f4616fa6e736ba37e9c219dbda637bfb5f89b))
* **lobby:** keep members up to date with lobby state ([ab12243](https://github.com/sceredi/co-type/commit/ab12243636786c8a5bef6ae376afb35a07d0c50a))
* register created lobbies to broker ([0b4a43e](https://github.com/sceredi/co-type/commit/0b4a43e4586f4c99fd63f1df1a4c8144d1dfe58f))

### Dependency updates

* **deps:** update alpine docker tag to v3.24 ([a9d6f35](https://github.com/sceredi/co-type/commit/a9d6f356c3c795ed852f2ecfc6f7a20dbd066211))
* **deps:** update module charm.land/bubbletea/v2 to v2.0.7 ([8b446f4](https://github.com/sceredi/co-type/commit/8b446f4b4094294f4d877c778aca0b667419b0e4))
* **deps:** update module charm.land/lipgloss/v2 to v2.0.4 ([b29b8fc](https://github.com/sceredi/co-type/commit/b29b8fc0a30782938a156789782f6e739b563397))
* **deps:** update node.js to 24.17 ([ddb3ced](https://github.com/sceredi/co-type/commit/ddb3ced21f76be8551f90a2d8cd8ecc949aea986))

### Bug Fixes

* avoid 2 players with the same name joining the same lobby ([713a7b3](https://github.com/sceredi/co-type/commit/713a7b3b4824fa1e9f95bedae11004567932af82))
* delete lobby if unable to register it on broker ([85e3477](https://github.com/sceredi/co-type/commit/85e3477a8153376d0bc6a0e279ccf881e33f53bf))

### Build and continuous integration

* **deps:** update actions/checkout action to v6.0.3 ([89ad0c0](https://github.com/sceredi/co-type/commit/89ad0c01705eebafe280aadf79e0a184aacb8b8c))
* **deps:** update actions/checkout action to v7 ([2a25549](https://github.com/sceredi/co-type/commit/2a25549985dc27cde375d8bae3a8b6839efc36fe))
* **deps:** update danysk/action-checkout action to v0.2.29 ([2d04191](https://github.com/sceredi/co-type/commit/2d04191e13f7168dd311e55741163299c0802b6f))

### General maintenance

* **lobby-creation:** simplify message values ([2e8ed1f](https://github.com/sceredi/co-type/commit/2e8ed1f0ec44f9e55ae4ca6e3943bdf81419e1b1))

## [0.3.0](https://github.com/sceredi/co-type/compare/v0.2.0...v0.3.0) (2026-06-02)

### Features

* add client broker discovery upon lobby creation ([ce43e71](https://github.com/sceredi/co-type/commit/ce43e710224bcab11611f34428bae708d23c7242))
* **client:** create and connect to lobby ([892a3df](https://github.com/sceredi/co-type/commit/892a3dfb26a83eb18f658023be150ff9196b549b))
* **client:** get best server from broker ([e4d2834](https://github.com/sceredi/co-type/commit/e4d28343e8c9a7bd27ee46d8584f390be0b9d5c9))
* containerize application ([30dd523](https://github.com/sceredi/co-type/commit/30dd523727f16e9824d870dc81bb862bb34add6a))
* **server:** add lobby creation ([e547ee6](https://github.com/sceredi/co-type/commit/e547ee622a8f346d7a8797e2be2de9984e14a0e9))
* **server:** split lobby creation and event subscription ([3f53b37](https://github.com/sceredi/co-type/commit/3f53b37305e4a2b163fff572b43e8205cee29da0))

### Dependency updates

* **deps:** update module google.golang.org/grpc to v1.81.0 ([eceff25](https://github.com/sceredi/co-type/commit/eceff259515b5f103efb5aed691bdb1f40d3dbd9))
* **deps:** update module google.golang.org/grpc to v1.81.1 ([d7853e7](https://github.com/sceredi/co-type/commit/d7853e79142840d097aa3c41e9f19a17b2721b18))
* **deps:** update node.js to 24.16 ([f661da8](https://github.com/sceredi/co-type/commit/f661da87b5cfcfd977f87e7ad8275add89cf2559))

### Bug Fixes

* correctly handle subscribe errors ([2c8fd8d](https://github.com/sceredi/co-type/commit/2c8fd8d6eec5b57cdb1025cb2421a930fed12016))

### Tests

* **client:** update gateway tests ([b84167c](https://github.com/sceredi/co-type/commit/b84167c4c6ced5cfc1f701bfb87e8bea3008f703))

### Build and continuous integration

* **deps:** update docker/build-push-action action to v7.2.0 ([a38d051](https://github.com/sceredi/co-type/commit/a38d0513efd2bff76bea8e5d3139420ed9d9c81b))
* **deps:** update docker/login-action action to v4.2.0 ([cab4e9b](https://github.com/sceredi/co-type/commit/cab4e9b20f47bbf1a0054d5ed13f0a355c769fc8))
* **deps:** update docker/setup-buildx-action action to v4.1.0 ([da188dc](https://github.com/sceredi/co-type/commit/da188dc2e67de7fa445c22cd6875da5e693985dc))
* **deps:** update golangci/golangci-lint-action action to v9.2.1 ([3dcff8c](https://github.com/sceredi/co-type/commit/3dcff8c77a9853acd360bf2cd1743fc0c4c51b40))

### General maintenance

* **broker:** move handlers to handler package ([30893dc](https://github.com/sceredi/co-type/commit/30893dc09e364fa0380976197155000b5d5d95e7))
* **client:** move gateway to its own package ([5be646b](https://github.com/sceredi/co-type/commit/5be646b42cb05785af11e15210b1bac18a809ae1))
* fix typo ([effefb8](https://github.com/sceredi/co-type/commit/effefb8905321bb2f05af78a9daa4432b1fc9422))
* move from k3s to kind ([ffd3482](https://github.com/sceredi/co-type/commit/ffd34823c5d10d1c5792859cc77fb7e7fa5079dc))

### Refactoring

* use int64 instead of 32 ([19b51d5](https://github.com/sceredi/co-type/commit/19b51d5565931619396ec998366fc0b8b734a2a9))

## [0.2.0](https://github.com/sceredi/co-type/compare/v0.1.0...v0.2.0) (2026-04-29)

### Features

* add name to servers ([5cc665f](https://github.com/sceredi/co-type/commit/5cc665f0d3b29f07fb9ba3f815e1f6ac2ad532c4))
* add server registration ([2fc09c4](https://github.com/sceredi/co-type/commit/2fc09c4f2d554e8114e7e9d182a16a8ea69a7208))
* **broker:** server creation ([1e033a1](https://github.com/sceredi/co-type/commit/1e033a1e9404548d643595ffa2e6a0225e35fe69))

### Bug Fixes

* manage client disconnenction ([c0bb4bc](https://github.com/sceredi/co-type/commit/c0bb4bc958bfc273113b24583c8de640d4980fd8))

### General maintenance

* fix typo ([a730c82](https://github.com/sceredi/co-type/commit/a730c82c582fa244c381597900950ccc13253ae8))
* fix typo ([dddd37f](https://github.com/sceredi/co-type/commit/dddd37fe5a127cdad87c8af14320de6b68c5f55d))
* group config utils ([6d1523d](https://github.com/sceredi/co-type/commit/6d1523d80de9e84a1ec5e57aedd750fa57ad7ff8))
* remove unreachable code ([1c072f4](https://github.com/sceredi/co-type/commit/1c072f4d6ed0ab109217b362d2372c98a25fe4be))

## [0.1.0](https://github.com/sceredi/co-type/compare/v0.0.0...v0.1.0) (2026-04-21)

### Features

* add end page ([0304706](https://github.com/sceredi/co-type/commit/0304706bb63cdb100a702732384bb960e86bd0a1))
* add keybind tooltip ([d189165](https://github.com/sceredi/co-type/commit/d189165299406c65ebf1347e25dab05eae565df5))
* add lobby-settings interaction ([f7dfb78](https://github.com/sceredi/co-type/commit/f7dfb78ac83fa5965fc0c17470868a48452248b1))
* add min term size check ([152cac6](https://github.com/sceredi/co-type/commit/152cac6c409d03187797190f3f2995226577118b))
* add paused game popup ([7b8be85](https://github.com/sceredi/co-type/commit/7b8be85766f0db5c4ef0515877bc84ddbb6da9ab))
* add player list to game screen ([dc4f272](https://github.com/sceredi/co-type/commit/dc4f2721ffda0f70ef63e357840f46b1c2e4d0bc))
* add player selection ([0bf8a02](https://github.com/sceredi/co-type/commit/0bf8a021380735074349d80aeb0061749fea24f1))
* add welcome-lobby pages movement ([bc64f18](https://github.com/sceredi/co-type/commit/bc64f18d271a47cc757840a02d672c5a51e13a4f))
* create new screen for settings editing ([d60c951](https://github.com/sceredi/co-type/commit/d60c95118d8a5d891fc04a6a756feb6b12d74ba4))
* create player list component ([29e3a34](https://github.com/sceredi/co-type/commit/29e3a34ae5ade1e62d4d0b67f923d00d89a8a010))
* create welcome page ([9588208](https://github.com/sceredi/co-type/commit/9588208e7a6d80b790a9f5fadd48d5da3860c0c1))
* mock game page typing section ([f55a23a](https://github.com/sceredi/co-type/commit/f55a23a717d24eefa75b451da163ea14b97059c2))
* mock welcome page ([42d1612](https://github.com/sceredi/co-type/commit/42d161215b61eb241eb493d839e529328f371dec))

### Dependency updates

* **deps:** update alpine docker tag to v3.23 ([081db89](https://github.com/sceredi/co-type/commit/081db89e9d59770a9cb7805b5f0cd949ec42165b))
* **deps:** update module charm.land/bubbletea/v2 to v2.0.6 ([0d13f78](https://github.com/sceredi/co-type/commit/0d13f782ced3185346be010352b5d911dc87cfaa))
* **deps:** update module charm.land/lipgloss/v2 to v2.0.3 ([fb38e4a](https://github.com/sceredi/co-type/commit/fb38e4aeab59945430f107b0d41ec434e3d46210))
* **deps:** update node.js to 24.15 ([f7bbde2](https://github.com/sceredi/co-type/commit/f7bbde2bb5a78c73502b608cecfadd0781865a5e))

### Bug Fixes

* correctly align lobby page ([ece2777](https://github.com/sceredi/co-type/commit/ece27774fb78c40fa9e0663dae908851a0066a81))
* divide by zero ([1b81061](https://github.com/sceredi/co-type/commit/1b8106193f5257a379fabe2e57e79034e688eb1f))
* package.json directives ([a9ed870](https://github.com/sceredi/co-type/commit/a9ed870c9e2d7b2d285896558a94023e4cd8d966))

### Tests

* add testing for lobby components ([72f52c9](https://github.com/sceredi/co-type/commit/72f52c9625912926e03637a2e43f231853ea4b74))
* add testing for welcome and main model ([1fe27da](https://github.com/sceredi/co-type/commit/1fe27da8a1a4fccacd4176f6aa8d2f4e89fd8934))

### Build and continuous integration

* always try to build website ([62c2bc4](https://github.com/sceredi/co-type/commit/62c2bc4263e445dfd17ca95c71457c921cba0208))
* base workflow ([d459c5a](https://github.com/sceredi/co-type/commit/d459c5a76ea09525c0973deee198b9fe9b4003d7))
* **deps:** update actions/checkout action to v4.3.1 ([cc1f37f](https://github.com/sceredi/co-type/commit/cc1f37f4ba9c8c5774879542e2d04f642281ca28))
* **deps:** update actions/checkout action to v5 ([4a163ff](https://github.com/sceredi/co-type/commit/4a163ff0b225d7ad45ade49b40d1bb06b75c446a))
* **deps:** update actions/checkout action to v6 ([6f55afa](https://github.com/sceredi/co-type/commit/6f55afa9c88e8f2d3dd5ae915f33c7361e271102))
* **deps:** update actions/setup-go action to v6.4.0 ([7fd54c8](https://github.com/sceredi/co-type/commit/7fd54c88600135bd6c73dd861b57e0421e0210ed))
* **deps:** update actions/setup-node action to v6.4.0 ([1625bcf](https://github.com/sceredi/co-type/commit/1625bcf127525e3cda7af1cb922dec02938f1fb2))
* **deps:** update docker/build-push-action action to v6.19.2 ([747e747](https://github.com/sceredi/co-type/commit/747e747fe5fe6328d36315ecaf6d0ae160414963))
* **deps:** update docker/build-push-action action to v7 ([8b2b283](https://github.com/sceredi/co-type/commit/8b2b2836f130c936c6ba2ec3e206b1721254cfb4))
* **deps:** update docker/build-push-action action to v7.1.0 ([2fb6f12](https://github.com/sceredi/co-type/commit/2fb6f1275725b784969c38aadd57cb5bbc06d723))
* **deps:** update docker/login-action action to v4.1.0 ([4492e26](https://github.com/sceredi/co-type/commit/4492e2601c77d56d50650558f77d3e0c9587fdb4))
* **deps:** update docker/metadata-action action to v6 ([16a7cc1](https://github.com/sceredi/co-type/commit/16a7cc1b5f522a78694f8389873ca7b925bc0d9b))
* **deps:** update docker/setup-buildx-action action to v4 ([71399de](https://github.com/sceredi/co-type/commit/71399de8736747e170b866e84657daffc4151247))
* update docker packages name ([0558618](https://github.com/sceredi/co-type/commit/05586189a54b95f104a451159086e39db8cb4a02))

### General maintenance

* add Dockerfile ([feb90b3](https://github.com/sceredi/co-type/commit/feb90b3cd7f3210bdbd0acb6a76bf2d72e4cd4d5))
* add golangci settings ([1c1c816](https://github.com/sceredi/co-type/commit/1c1c816fb224470d517b03f24fd51317e3ae425b))
* add initial text to settings ([32f7ac6](https://github.com/sceredi/co-type/commit/32f7ac6598cb7f81055a19b29fa9218db02672b1))
* add linter configuration and fix wanings ([aa111f0](https://github.com/sceredi/co-type/commit/aa111f01bdfafcd23a3b2f0c620d397d37f82925))
* add writing to settings inputs ([1010b20](https://github.com/sceredi/co-type/commit/1010b20c4a0e6e0f747446da9babb24d2377f487))
* center interface ([76a12c3](https://github.com/sceredi/co-type/commit/76a12c352cb5179ece593ce164c3359927d6e2f7))
* fix linter errors ([e3aebe5](https://github.com/sceredi/co-type/commit/e3aebe53ac3abf917910a786a6e44d850d9f38d6))
* fix typo ([1b185c8](https://github.com/sceredi/co-type/commit/1b185c8e233707c6600ab640f2e1992e7397cc2e))
* make header stateless ([9d120b7](https://github.com/sceredi/co-type/commit/9d120b77a00b4cbcbdbcc001b0f70a3a188ffbd4))
* mock header and ready/leave buttons ([fe7732e](https://github.com/sceredi/co-type/commit/fe7732e9cd97c05677528f932af3630d18dd4849))
* remove example code ([60a758b](https://github.com/sceredi/co-type/commit/60a758b231a840d14b6caad9c7eef47655579cfc))
* remove tea.Cursor setup ([9b37bb6](https://github.com/sceredi/co-type/commit/9b37bb68d4fb21345c46581f03f13d30a317ce77))
* style ready button ([79b9fa7](https://github.com/sceredi/co-type/commit/79b9fa7275695ddfe0c979255b5f1d83cb968c69))
* use player pointers in lobby ([dbdb7f9](https://github.com/sceredi/co-type/commit/dbdb7f9243e56f0929e7008f28e6e549a330001c))
