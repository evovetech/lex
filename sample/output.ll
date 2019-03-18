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

define double @min(double %x, double %y) {
entry:
  %cmptmp = fcmp ult double %x, %y
  %x.y = select i1 %cmptmp, double %x, double %y
  ret double %x.y
}

define double @max(double %x, double %y) {
entry:
  %cmptmp = fcmp ult double %x, %y
  %y.x = select i1 %cmptmp, double %y, double %x
  ret double %y.x
}

define double @fib(double %x) {
entry:
  %cmptmp = fcmp ult double %x, 3.000000e+00
  br i1 %cmptmp, label %ifcont, label %else

else:                                             ; preds = %entry
  %subtmp = fadd double %x, -1.000000e+00
  %calltmp = call double @fib(double %subtmp)
  %subtmp1 = fadd double %x, -2.000000e+00
  %calltmp2 = call double @fib(double %subtmp1)
  %addtmp = fadd double %calltmp, %calltmp2
  br label %ifcont

ifcont:                                           ; preds = %else, %entry
  %iftmp = phi double [ %addtmp, %else ], [ 1.000000e+00, %entry ]
  ret double %iftmp
}
