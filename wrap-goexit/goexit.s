// Copyright 2018 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

#include "go_asm.h"
#include "textflag.h"

// pc和入口相等
TEXT ·hackedGoexit(SB),NOSPLIT,$0-0
	CALL	·hackedGoexit1(SB)
