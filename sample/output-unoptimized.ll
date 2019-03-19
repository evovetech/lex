; ModuleID = 'output-unoptimized.bc'
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
  %booltmp = uitofp i1 %cmptmp to double
  %ifcond = fcmp one double %booltmp, 0.000000e+00
  br i1 %ifcond, label %then, label %else

then:                                             ; preds = %entry
  br label %ifcont

else:                                             ; preds = %entry
  br label %ifcont

ifcont:                                           ; preds = %else, %then
  %iftmp = phi double [ %x, %then ], [ %y, %else ]
  ret double %iftmp
}

define double @max(double %x, double %y) {
entry:
  %cmptmp = fcmp ult double %x, %y
  %booltmp = uitofp i1 %cmptmp to double
  %ifcond = fcmp one double %booltmp, 0.000000e+00
  br i1 %ifcond, label %then, label %else

then:                                             ; preds = %entry
  br label %ifcont

else:                                             ; preds = %entry
  br label %ifcont

ifcont:                                           ; preds = %else, %then
  %iftmp = phi double [ %y, %then ], [ %x, %else ]
  ret double %iftmp
}

define double @fib(double %x) {
entry:
  %cmptmp = fcmp ult double %x, 3.000000e+00
  %booltmp = uitofp i1 %cmptmp to double
  %ifcond = fcmp one double %booltmp, 0.000000e+00
  br i1 %ifcond, label %then, label %else

then:                                             ; preds = %entry
  br label %ifcont

else:                                             ; preds = %entry
  %subtmp = fsub double %x, 1.000000e+00
  %calltmp = call double @fib(double %subtmp)
  %subtmp1 = fsub double %x, 2.000000e+00
  %calltmp2 = call double @fib(double %subtmp1)
  %addtmp = fadd double %calltmp, %calltmp2
  br label %ifcont

ifcont:                                           ; preds = %else, %then
  %iftmp = phi double [ 1.000000e+00, %then ], [ %addtmp, %else ]
  ret double %iftmp
}

declare double @sin(double)

declare double @cos(double)

declare double @atan2(double, double)

define double @callAtan2(double %x, double %y) {
entry:
  %calltmp = call double @sin(double %x)
  %calltmp1 = call double @cos(double %y)
  %calltmp2 = call double @atan2(double %calltmp, double %calltmp1)
  ret double %calltmp2
}
