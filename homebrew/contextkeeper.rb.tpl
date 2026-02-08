# Homebrew formula template for ContextKeeper
# Copy to your tap as contextkeeper.rb and update SHA256

class Contextkeeper < Formula
  desc "A minimalist CLI tool for managing project context"
  homepage "https://github.com/ondrahracek/contextkeeper"
  
  # Update URL and SHA256 for each release
  url "https://github.com/ondrahracek/contextkeeper/releases/download/v#{VERSION}/contextkeeper-#{VERSION}-darwin-arm64.tar.gz"
  sha256 "TODO: Run `sha256sum filename.tar.gz` and add here"
  license "MIT"
  version "#{VERSION}"

  def install
    bin.install "contextkeeper"
  end

  test do
    system "#{bin}/contextkeeper", "--help"
  end
end
