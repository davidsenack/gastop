# Homebrew Formula for gastop
# To install locally: brew install --build-from-source ./packaging/homebrew/gastop.rb
# To tap: brew tap davidsenack/gastop && brew install gastop

class Gastop < Formula
  desc "htop-like terminal UI for Gas Town workspaces"
  homepage "https://github.com/davidsenack/gastop"
  url "https://github.com/davidsenack/gastop/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "7cef48c21801f8c11a7e5fb071700998ac83032651799a1be5e60a8f78598c1e"
  license "MIT"
  head "https://github.com/davidsenack/gastop.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X main.version=#{version}"), "./cmd/gastop"
  end

  test do
    assert_match "gastop", shell_output("#{bin}/gastop --help 2>&1", 0)
  end
end
