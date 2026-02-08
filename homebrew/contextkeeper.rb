# Homebrew formula for ContextKeeper
# A minimalist CLI tool for managing project context

class Contextkeeper < Formula
  desc "A minimalist CLI tool for managing project context"
  homepage "https://github.com/ondrahracek/contextkeeper"
  url "https://github.com/ondrahracek/contextkeeper/releases/download/v0.2.0/contextkeeper-0.2.0-src.tar.gz"
  sha256 "TODO: Add SHA256 after first release"
  license "MIT"
  version "0.2.0"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", "ck", "."
    bin.install "ck"
  end

  test do
    system "#{bin}/ck", "--help"
  end
end
