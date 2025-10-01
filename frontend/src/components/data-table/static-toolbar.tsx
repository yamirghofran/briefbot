import * as React from "react";
import { Search, Filter, Zap, Check, Loader2 } from "lucide-react";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { UrlSubmissionDialog } from "@/components/url-submission-dialog";
import { digestApi } from "@/services/api";

interface StaticToolbarProps extends React.HTMLAttributes<HTMLDivElement> {
  userId?: number;
}

export function StaticToolbar({
  className,
  userId,
  ...props
}: StaticToolbarProps) {
  const [searchValue, setSearchValue] = React.useState("");
  const [showSuccess, setShowSuccess] = React.useState(false);

  // Mutation for triggering integrated digest
  const triggerDigestMutation = useMutation({
    mutationFn: () => {
      if (userId) {
        return digestApi.triggerIntegratedDigestForUser(userId);
      } else {
        return digestApi.triggerIntegratedDigest();
      }
    },
    onSuccess: () => {
      toast.success(
        "Digest processing started! You'll receive an email when ready.",
      );
      setShowSuccess(true);
      // Reset success state after 3 seconds
      setTimeout(() => {
        setShowSuccess(false);
      }, 3000);
    },
    onError: (error) => {
      toast.error("Failed to trigger integrated digest");
      console.error("Digest trigger error:", error);
    },
  });

  const handleTriggerDigest = () => {
    triggerDigestMutation.mutate();
  };

  return (
    <div className="flex items-center justify-between" {...props}>
      <div className="flex flex-1 items-center space-x-2"></div>

      <div className="flex items-center space-x-2">
        {userId && <UrlSubmissionDialog userId={userId} />}
        <Button
          onClick={handleTriggerDigest}
          disabled={triggerDigestMutation.isPending || showSuccess}
          variant="default"
          size="sm"
          className={`transition-all duration-300 ${
            showSuccess
              ? "bg-green-600 hover:bg-green-600"
              : "bg-black hover:bg-gray-800"
          } text-white`}
        >
          {triggerDigestMutation.isPending ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Creating Digest...
            </>
          ) : showSuccess ? (
            <>
              <Check className="mr-2 h-4 w-4" />
              Digest Started!
            </>
          ) : (
            <>
              <Zap className="mr-2 h-4 w-4" />
              Trigger Digest
            </>
          )}
        </Button>
      </div>
    </div>
  );
}
