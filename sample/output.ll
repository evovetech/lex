; ModuleID = 'output.bc'
source_filename = "lex"
target datalayout = "e-m:o-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-apple-darwin17.7.0"

define double @sum(double %x, double %y) {
entry:
  %addtmp = fadd double %x, %y
  ret double %addtmp
}

define double @average(double %x, double %y) {
entry:
  %calltmp = call double @sum(double %x, double %y)
  %multmp = fmul double %calltmp, 5.000000e-01
  ret double %multmp
}
